package sprite

import(
	"time"
	"bytes"
	"image"
	"image/png"
	"github.com/disintegration/imaging"
)

//Types
type Cycle struct {
	Name string
	Index int
	Speed int
	Reverse bool
	Loop bool
	Mirror bool
	
	//pvt
	sprite *Sprite
	frames [][]byte
	mirrors [][]byte
	total int
	playing bool
	stopPlaying chan bool
}

//Create (frames are 0 indexed!)
func NewCycle(name string, sprite *Sprite, startFrame int, endFrame int, speed int) *Cycle {
	cycle := new(Cycle)
	
	cycle.Name = name
	cycle.Index = 0
	cycle.frames = cycle.getCycleFrames(sprite.Frames, startFrame, endFrame)
	cycle.mirrors = cycle.getCycleMirrors()
	cycle.total = len(cycle.frames)
	cycle.Speed = speed
	cycle.Loop = true
	cycle.Reverse = false
	cycle.sprite = sprite
	cycle.stopPlaying  = make(chan bool, 10)
	
	return cycle
}

func AllFramesCycle(sprite *Sprite) *Cycle {
	cycle := new(Cycle)
	
	cycle.Name = "All"
	cycle.Index = 0
	cycle.frames = cycle.getCycleFrames(sprite.Frames, 0, len(sprite.Frames))
	cycle.total = len(cycle.frames)
	cycle.Speed = 1
	cycle.Loop = true
	cycle.Reverse = false
	cycle.sprite = sprite
	cycle.stopPlaying  = make(chan bool, 10)
	
	return cycle
}

//Utils
func (c *Cycle) getCycleFrames(frames [][]byte, start int, end int) [][]byte {
	var cycleFrames [][]byte
	for i := start; i < end; i++ {
		cycleFrames = append(cycleFrames, frames[i])
	}
	return cycleFrames
}

func (c *Cycle) getCycleMirrors() [][]byte {
	var mirrors [][]byte
	for _, frame := range c.frames {
		newMirror := c.getMirror(frame)
		mirrors = append(mirrors, newMirror)
	}
	return mirrors
}

func (c *Cycle) getMirror(frame []byte) []byte {
	newImage, _, _ := image.Decode(bytes.NewReader(frame))
	flipped := imaging.FlipH(newImage)
	
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, flipped)
	return buf.Bytes()
}

func (c *Cycle) RefreshFrame() {
	if c.Mirror {
		c.sprite.SetImage(c.mirrors[c.Index])
	} else {
		c.sprite.SetImage(c.frames[c.Index])
	}
}

func (c *Cycle) IsPlaying() bool {
	return c.playing
}


//Actions
func (c *Cycle) Play() {
	c.playing = true
	go func() {
		for {
			select {
				case <- c.stopPlaying:
					return
				default:
					if c.Reverse {
						c.Prev()
					} else {
						c.Next()
					}
					time.Sleep(time.Second / time.Duration(c.Speed))
	   		 }
		}
	}()
}

func (c *Cycle) Stop() {
	c.stopPlaying <- true
	c.playing = false
	c.Index = 0
	c.RefreshFrame()
}

func (c *Cycle) Next() {
	if c.Index < c.total - 1 {
		c.Index++
	} else {
		//If looping continue loop, else stop playing
		if c.Loop {
			c.Index = 0
		}
	}
	c.RefreshFrame()
}

func (c *Cycle) Prev() {
	//If looping continue loop, else stop playing
	if c.Index > 0 {
		c.Index--
	} else {
		if c.Loop {
			c.Index = c.total - 1
		}
	}
	c.RefreshFrame()
}