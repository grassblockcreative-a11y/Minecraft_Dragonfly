package generator

import (
	"math"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/biome"
	"github.com/df-mc/dragonfly/server/world/chunk"
)

type chunkRand struct {
	state uint64
}

func newChunkRand(seed uint64) *chunkRand {
	return &chunkRand{state: seed}
}

func (r *chunkRand) Uint32() uint32 {
	r.state = r.state*6364136223846793005 + 1442695040888963407
	return uint32(r.state >> 32)
}

func (r *chunkRand) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	return int(r.Uint32() % uint32(n))
}

func chunkSeed(cx, cz int32) uint64 {
	return uint64(cx)*341873128712 + uint64(cz)*132897987541
}

// VanillaLightweight is a lightweight vanilla-like terrain generator optimised
// for low-end ARM boards (e.g. Armbian H680P).
type VanillaLightweight struct {
	biome     uint32
	grass     uint32
	dirt      uint32
	stone     uint32
	sand      uint32
	water     uint32
	bedrock   uint32
	deepslate uint32

	// Ores (Stone variants)
	coalOre     uint32
	ironOre     uint32
	goldOre     uint32
	diamondOre  uint32
	lapisOre    uint32
	redstoneOre uint32
	copperOre   uint32
	emeraldOre  uint32

	// Ores (Deepslate variants)
	coalOreDeepslate     uint32
	ironOreDeepslate     uint32
	goldOreDeepslate     uint32
	diamondOreDeepslate  uint32
	lapisOreDeepslate    uint32
	redstoneOreDeepslate uint32
	copperOreDeepslate   uint32
	emeraldOreDeepslate  uint32

	// Vegetation & Trees
	shortGrass uint32
	fern       uint32
	dandelion  uint32
	log        uint32
	leaves     uint32

	// Village Houses
	cobblestone uint32
	planks      uint32
	glass       uint32
}

// NewVanillaLightweight creates a new VanillaLightweight generator.
func NewVanillaLightweight(br world.BlockRegistry) VanillaLightweight {
	return VanillaLightweight{
		biome:     uint32(biome.Plains{}.EncodeBiome()),
		grass:     br.BlockRuntimeID(block.Grass{}),
		dirt:      br.BlockRuntimeID(block.Dirt{}),
		stone:     br.BlockRuntimeID(block.Stone{}),
		sand:      br.BlockRuntimeID(block.Sand{}),
		water:     br.BlockRuntimeID(block.Water{Still: true, Depth: 8}),
		bedrock:   br.BlockRuntimeID(block.Bedrock{}),
		deepslate: br.BlockRuntimeID(block.Deepslate{Type: block.NormalDeepslate()}),

		coalOre:     br.BlockRuntimeID(block.CoalOre{Type: block.StoneOre()}),
		ironOre:     br.BlockRuntimeID(block.IronOre{Type: block.StoneOre()}),
		goldOre:     br.BlockRuntimeID(block.GoldOre{Type: block.StoneOre()}),
		diamondOre:  br.BlockRuntimeID(block.DiamondOre{Type: block.StoneOre()}),
		lapisOre:    br.BlockRuntimeID(block.LapisOre{Type: block.StoneOre()}),
		redstoneOre: br.BlockRuntimeID(block.RedstoneOre{Type: block.StoneOre()}),
		copperOre:   br.BlockRuntimeID(block.CopperOre{Type: block.StoneOre()}),
		emeraldOre:  br.BlockRuntimeID(block.EmeraldOre{Type: block.StoneOre()}),

		coalOreDeepslate:     br.BlockRuntimeID(block.CoalOre{Type: block.DeepslateOre()}),
		ironOreDeepslate:     br.BlockRuntimeID(block.IronOre{Type: block.DeepslateOre()}),
		goldOreDeepslate:     br.BlockRuntimeID(block.GoldOre{Type: block.DeepslateOre()}),
		diamondOreDeepslate:  br.BlockRuntimeID(block.DiamondOre{Type: block.DeepslateOre()}),
		lapisOreDeepslate:    br.BlockRuntimeID(block.LapisOre{Type: block.DeepslateOre()}),
		redstoneOreDeepslate: br.BlockRuntimeID(block.RedstoneOre{Type: block.DeepslateOre()}),
		copperOreDeepslate:   br.BlockRuntimeID(block.CopperOre{Type: block.DeepslateOre()}),
		emeraldOreDeepslate:  br.BlockRuntimeID(block.EmeraldOre{Type: block.DeepslateOre()}),

		shortGrass: br.BlockRuntimeID(block.ShortGrass{}),
		fern:       br.BlockRuntimeID(block.Fern{}),
		dandelion:  br.BlockRuntimeID(block.Flower{Type: block.Dandelion()}),
		log:        br.BlockRuntimeID(block.Log{Wood: block.OakWood()}),
		leaves:     br.BlockRuntimeID(block.Leaves{Type: block.OakLeaves()}),

		cobblestone: br.BlockRuntimeID(block.Cobblestone{}),
		planks:      br.BlockRuntimeID(block.Planks{Wood: block.OakWood()}),
		glass:       br.BlockRuntimeID(block.Glass{}),
	}
}

const seaLevel int16 = 62

// heightAt returns the terrain height at a given world coordinate using cheap
// trigonometric functions. Supports mountains.
func heightAt(wx, wz float64) int16 {
	h := 64.0 +
		8.0*math.Sin(wx*0.0157)*math.Cos(wz*0.0157) +
		4.0*math.Sin(wx*0.0341+wz*0.0271) +
		2.0*math.Cos(wx*0.083)*math.Sin(wz*0.083)

	// Mountains (low frequency wave)
	mountainNoise := math.Sin(wx*0.005) * math.Cos(wz*0.005)
	if mountainNoise > 0.1 {
		h += (mountainNoise - 0.1) * 60.0
	}
	return int16(h)
}

// isCave determines if there should be a cave at the given absolute coordinate.
func isCave(wx, wy, wz float64, height, min int16) bool {
	if int16(wy) <= min+4 || int16(wy) >= height-4 {
		return false
	}
	// Fast trigonometric cave formula
	n := math.Sin(wx*0.09)*math.Cos(wz*0.09)*math.Sin(wy*0.15) +
		math.Cos(wx*0.04)*math.Sin(wz*0.04)*math.Cos(wy*0.08)
	return n > 0.72
}

// getOreRuntimeID determines if a block should be an ore.
func (v VanillaLightweight) getOreRuntimeID(wx float64, y int16, wz float64) (uint32, bool) {
	h := uint32(wx)*374761393 + uint32(y)*668265263 + uint32(wz)*131071
	h = (h ^ (h >> 13)) * 12741261
	h = h & 0xffff

	isDeep := y < 0

	if y < 16 {
		if h < 8 {
			if isDeep {
				return v.diamondOreDeepslate, true
			}
			return v.diamondOre, true
		}
	}
	if y < 32 {
		if h < 20 {
			if isDeep {
				return v.goldOreDeepslate, true
			}
			return v.goldOre, true
		}
	}
	if y < 32 {
		if h < 18 {
			if isDeep {
				return v.lapisOreDeepslate, true
			}
			return v.lapisOre, true
		}
	}
	if y < 16 {
		if h < 30 {
			if isDeep {
				return v.redstoneOreDeepslate, true
			}
			return v.redstoneOre, true
		}
	}
	if y < 64 {
		if h < 55 {
			if isDeep {
				return v.ironOreDeepslate, true
			}
			return v.ironOre, true
		}
	}
	if y < 48 {
		if h < 45 {
			if isDeep {
				return v.copperOreDeepslate, true
			}
			return v.copperOre, true
		}
	}
	if y < 128 {
		if h < 100 {
			if isDeep {
				return v.coalOreDeepslate, true
			}
			return v.coalOre, true
		}
	}
	if y > 60 && y < 120 {
		if h < 15 {
			return v.emeraldOre, true
		}
	}

	return 0, false
}

// GenerateChunk generates a single chunk.
func (v VanillaLightweight) GenerateChunk(pos world.ChunkPos, c *chunk.Chunk) {
	min := int16(c.Range().Min())
	cx, cz := int32(pos.X())<<4, int32(pos.Z())<<4

	minX := int32(pos.X()) * 16
	maxX := minX + 15
	minZ := int32(pos.Z()) * 16
	maxZ := minZ + 15

	// 1. Generate base terrain and caves
	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {
			wx := float64(cx + int32(x))
			wz := float64(cz + int32(z))

			height := heightAt(wx, wz)

			top := height
			if seaLevel > top {
				top = seaLevel
			}

			for y := min; y <= top; y++ {
				c.SetBiome(x, y, z, v.biome)

				if y == min {
					c.SetBlock(x, y, z, 0, v.bedrock)
					continue
				}

				// Check for cave carving
				if isCave(wx, float64(y), wz, height, min) {
					c.SetBlock(x, y, z, 0, 0) // Air
					continue
				}

				switch {
				case y < height-3:
					if oreID, ok := v.getOreRuntimeID(wx, y, wz); ok {
						c.SetBlock(x, y, z, 0, oreID)
					} else if y < 0 {
						c.SetBlock(x, y, z, 0, v.deepslate)
					} else {
						c.SetBlock(x, y, z, 0, v.stone)
					}
				case y < height:
					c.SetBlock(x, y, z, 0, v.dirt)
				case y == height:
					if height < seaLevel-1 {
						// Underwater floor → sand.
						c.SetBlock(x, y, z, 0, v.sand)
					} else {
						c.SetBlock(x, y, z, 0, v.grass)

						// Spawn short grass / fern / dandelion with a coordinate-based hash
						hGrass := uint32(wx)*131071 + uint32(wz)*374761393
						hGrass = (hGrass ^ (hGrass >> 13)) * 12741261
						hGrass = hGrass & 0xff
						if hGrass < 40 && height+1 <= int16(c.Range().Max()) {
							if hGrass < 30 {
								c.SetBlock(x, height+1, z, 0, v.shortGrass)
							} else if hGrass < 35 {
								c.SetBlock(x, height+1, z, 0, v.fern)
							} else {
								c.SetBlock(x, height+1, z, 0, v.dandelion)
							}
						}
					}
				default:
					// y > height && y <= seaLevel → water.
					if y <= seaLevel {
						c.SetBlock(x, y, z, 0, v.water)
					}
				}
			}

			// Set biome above terrain
			for y := top + 1; y <= top+16 && y <= int16(c.Range().Max()); y++ {
				c.SetBiome(x, y, z, v.biome)
			}
		}
	}

	// 2. Generate structures (trees, village houses) using neighbor lookup for seamless chunk boundaries
	for ncx := pos.X() - 1; ncx <= pos.X() + 1; ncx++ {
		for ncz := pos.Z() - 1; ncz <= pos.Z() + 1; ncz++ {
			seed := chunkSeed(int32(ncx), int32(ncz))
			rng := newChunkRand(seed)

			// 2a. Village house (3% chance per chunk)
			hasHouse := rng.Intn(100) < 3
			if hasHouse {
				hx := int32(ncx)*16 + int32(rng.Intn(10)) + 3
				hz := int32(ncz)*16 + int32(rng.Intn(10)) + 3
				hy := heightAt(float64(hx), float64(hz))

				if hy >= seaLevel {
					// Draw the house
					for dx := int32(-2); dx <= 2; dx++ {
						for dz := int32(-2); dz <= 2; dz++ {
							ax := hx + dx
							az := hz + dz
							if ax >= minX && ax <= maxX && az >= minZ && az <= maxZ {
								lx := uint8(ax - minX)
								lz := uint8(az - minZ)
								for dy := int16(0); dy <= 4; dy++ {
									ay := hy + dy
									var bid uint32
									if dx == -2 || dx == 2 || dz == -2 || dz == 2 {
										// Walls
										if dy == 0 {
											bid = v.cobblestone
										} else if dy < 4 {
											if (dx == -2 || dx == 2) && (dz == -2 || dz == 2) {
												bid = v.cobblestone
											} else if dy == 2 && (dx == 0 || dz == 0) {
												bid = v.glass
											} else if dy == 1 && dx == 0 && dz == -2 {
												bid = 0 // Door bottom
											} else if dy == 2 && dx == 0 && dz == -2 {
												bid = 0 // Door top
											} else {
												bid = v.planks
											}
										} else {
											bid = v.planks
										}
									} else {
										bid = 0 // Inside air
									}
									if ay >= min && ay <= int16(c.Range().Max()) {
										c.SetBlock(lx, ay, lz, 0, bid)
									}
								}
							}
						}
					}
				}
				// Skip generating trees in chunks that have houses to avoid overlap
				continue
			}

			// 2b. Trees (0-2 per chunk)
			numTrees := rng.Intn(3)
			for i := 0; i < numTrees; i++ {
				tx := int32(ncx)*16 + int32(rng.Intn(16))
				tz := int32(ncz)*16 + int32(rng.Intn(16))
				ty := heightAt(float64(tx), float64(tz))

				if ty >= seaLevel {
					th := int16(4 + rng.Intn(3)) // height 4-6

					// Draw Trunk
					if tx >= minX && tx <= maxX && tz >= minZ && tz <= maxZ {
						lx := uint8(tx - minX)
						lz := uint8(tz - minZ)
						for dy := int16(1); dy <= th; dy++ {
							ay := ty + dy
							if ay >= min && ay <= int16(c.Range().Max()) {
								c.SetBlock(lx, ay, lz, 0, v.log)
							}
						}
					}

					// Draw Leaves
					for ly := ty + th - 2; ly <= ty + th + 1; ly++ {
						r := int32(2)
						if ly > ty+th {
							r = 1
						}
						for dx := -r; dx <= r; dx++ {
							for dz := -r; dz <= r; dz++ {
								if r == 2 && (dx == -2 || dx == 2) && (dz == -2 || dz == 2) {
									continue // Round leaves corners
								}
								if dx == 0 && dz == 0 && ly <= ty+th {
									continue // Trunk position
								}
								ax := tx + dx
								az := tz + dz
								if ax >= minX && ax <= maxX && az >= minZ && az <= maxZ {
									lx := uint8(ax - minX)
									lz := uint8(az - minZ)
									if ly >= min && ly <= int16(c.Range().Max()) {
										c.SetBlock(lx, ly, lz, 0, v.leaves)
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

// DefaultSpawn returns a fixed, safe spawn position on solid ground at 0, 0.
func (v VanillaLightweight) DefaultSpawn(dim world.Dimension) cube.Pos {
	h := heightAt(0, 0)
	return cube.Pos{0, int(h) + 1, 0}
}
