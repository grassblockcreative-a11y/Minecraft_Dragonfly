package main

import (
	"fmt"
	"log/slog"
	"math"
	"math/rand/v2"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/pelletier/go-toml"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	conf, err := readConfig(slog.Default())
	if err != nil {
		panic(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()

	srv.Listen()

	// 1. Periodic mob spawner loop
	go func() {
		for {
			time.Sleep(20 * time.Second)
			srv.World().Exec(func(tx *world.Tx) {
				mobsCount := 0
				for e := range tx.Entities() {
					if _, ok := e.H().Type().(entity.MobType); ok {
						mobsCount++
					}
				}

				// Only spawn if count is below 25 to optimize performance
				if mobsCount < 25 {
					for plEntity := range tx.Players() {
						pl, ok := plEntity.(*player.Player)
						if !ok {
							continue
						}
						// 50% chance to spawn a mob around this player
						if rand.IntN(2) == 0 {
							dx := float64(rand.IntN(30) - 15)
							dz := float64(rand.IntN(30) - 15)
							if math.Abs(dx) > 5 || math.Abs(dz) > 5 {
								pos := pl.Position().Add(mgl64.Vec3{dx, 0, dz})
								y := float64(tx.HighestBlock(int(pos.X()), int(pos.Z()))) + 1.0
								if y > float64(tx.Range()[0]) {
									spawnPos := mgl64.Vec3{pos.X(), y, pos.Z()}
									types := []world.EntityType{
										entity.PigType,
										entity.CowType,
										entity.SheepType,
										entity.ZombieType,
										entity.SkeletonType,
									}
									t := types[rand.IntN(len(types))]

									opts := world.EntitySpawnOpts{Position: spawnPos}
									h := opts.New(t, entity.PassiveBehaviourConfig{
										Gravity: 0.08,
										Drag:    0.02,
									})
									tx.AddEntity(h)
								}
							}
						}
					}
				}
			})
		}
	}()

	// 2. Periodic item clearing and auto-optimization loop
	go func() {
		for {
			// Sleep 10 minutes, then notify
			time.Sleep(10 * time.Minute)
			srv.World().Exec(func(tx *world.Tx) {
				for plEntity := range tx.Players() {
					if pl, ok := plEntity.(*player.Player); ok {
						pl.Message("§e[Pemberitahuan] Semua item di tanah akan dibersihkan dalam 5 menit!§r")
					}
				}
			})

			// Sleep 5 minutes, then clear and optimize
			time.Sleep(5 * time.Minute)
			srv.World().Exec(func(tx *world.Tx) {
				count := 0
				for e := range tx.Entities() {
					if e.H().Type().EncodeEntity() == "minecraft:item" {
						_ = e.Close()
						count++
					}
				}
				for plEntity := range tx.Players() {
					if pl, ok := plEntity.(*player.Player); ok {
						pl.Message(fmt.Sprintf("§a[Pemberitahuan] Berhasil membersihkan %d item di tanah!§r", count))
					}
				}
			})

			// Trigger automatic optimization/GC
			runtime.GC()
			debug.FreeOSMemory()
		}
	}()

	for p := range srv.Accept() {
		p.SetGameMode(world.GameModeSurvival)
		p.Message("§aSelamat datang! Bermainlah dalam Survival Mode.§r")

		// Teleport player randomly near an existing player, or spawn
		go func(pl *player.Player) {
			time.Sleep(100 * time.Millisecond)
			pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {
				p := e.(*player.Player)
				var players []world.Entity
				for other := range tx.Players() {
					if other.H().UUID() != p.H().UUID() {
						players = append(players, other)
					}
				}

				var basePos mgl64.Vec3
				if len(players) > 0 {
					target := players[rand.IntN(len(players))]
					basePos = target.Position()
				} else {
					basePos = tx.World().Spawn().Vec3()
				}

				// Spawn randomly within 40 blocks of base position
				dx := float64(rand.IntN(80) - 40)
				dz := float64(rand.IntN(80) - 40)
				spawnX := basePos.X() + dx
				spawnZ := basePos.Z() + dz
				spawnY := float64(tx.HighestBlock(int(spawnX), int(spawnZ))) + 1.0

				spawnPos := mgl64.Vec3{spawnX, spawnY, spawnZ}
				p.Teleport(spawnPos)
			})
		}(p)
	}
}

func readConfig(log *slog.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config

	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}

		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}

		return c.Config(log)
	}

	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}

	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}

	return c.Config(log)
}
