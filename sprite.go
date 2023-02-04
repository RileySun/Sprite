package sprite

import(
	"fmt"
	"bytes"
	"image"
	"image/png"
	
	"github.com/oliamb/cutter"
)

type Sprite struct {
	Image []byte		//Current Sprite Frame
	Sheet []byte		//Sprite Sheet
	Frames [][]byte		//Frames of Sprite Sheet
	Cycle *Cycle		//Active Frame Cycle
	Cycles []*Cycle		//List of Frame Cycles
	Total int			//Number of Frames
	Rows int			//Number of Frame Rows (Only public because Cols is)
	Cols int			//Number of Frame Columns (may need this for grid display)
	OnUpdate func()		//Triggers anytime Image field is updated
	
	//pvt
	width int			//Width of Frame
	height int			//Height of Frame
}

//Create
func NewSprite(spriteSheet []byte, frameTotal int, rows int, columns int, frameWidth int, frameHeight int) *Sprite {
	sprite := new(Sprite)
	
	//Set Data
	sprite.Sheet = spriteSheet
	sprite.Total = frameTotal
	sprite.Rows = rows
	sprite.Cols = columns
	sprite.width = frameWidth
	sprite.height = frameHeight
	//Get Frames
	sprite.Frames = sprite.getFrames()
	sprite.SetImage(sprite.Frames[0])
	sprite.AddCycle(AllFramesCycle(sprite))//Default All Frames Cycle
	sprite.SetCycle("All")//Load Default Cycle
	
	return sprite
}

//Utils
func (s *Sprite) getFrames() [][]byte {
	var frames [][]byte
	
	reader := bytes.NewReader(s.Sheet)
	img, _, _ := image.Decode(reader) //image, ext, err
	
	var total, i, j int = 0, 0, 0
	
	for i = 0; i < s.Rows; i++ {
		//In case not all tiles are filled in
		if total > s.Total {
			return frames
		}
		
		for j = 0; j < s.Cols; j++ {
			cropped, _ := cutter.Crop(img, cutter.Config{
				Width: s.width,
				Height: s.height,
				Anchor: image.Point{j * s.width, i * s.height},
			})
			
			buffer := new(bytes.Buffer)
			_ = png.Encode(buffer, cropped)
			byt := buffer.Bytes()
			frames = append(frames, byt)
			total++
		}
	}
	
	return frames
}

func (s *Sprite) getCycle(name string) *Cycle {
	for _, cycle := range s.Cycles {
		if cycle.Name == name {
			return cycle
		}
	}
	return nil
}

//Actions
func (s *Sprite) Refresh() {
	if s.Cycle == nil {
		return
	} else {
		s.Cycle.RefreshFrame()
	}
} //Force refresh, useful in some situations

func (s *Sprite) SetImage(newImage []byte) {
	s.Image = newImage
	if s.OnUpdate != nil {
		s.OnUpdate()
	}
}

func (s *Sprite) ListCycles() []string {
	var list []string
	for _, cycle := range s.Cycles {
		list = append(list, cycle.Name)
	}
	return list
}

func (s *Sprite) AddCycle(cycle *Cycle) {
	test := s.getCycle(cycle.Name)
	
	if test != nil {
		fmt.Println("sprite.go - Duplicate Cycle Attempt, Can Not Re-Use Cycle Names")
	} else {
		s.Cycles = append(s.Cycles, cycle)
	}
}

func (s *Sprite) SetCycle(name string) {
	newCycle := s.getCycle(name)
	
	if newCycle == nil {
		fmt.Println("sprite.go - No such cycle exists, add cycle to use")
	} else {
		s.Cycle = newCycle
		s.Image = s.Cycle.frames[0]
	}
}

func (s *Sprite) Play() {
	if s.Cycle == nil {
		return
	} else {
		s.Cycle.Play()
	}
}

func (s *Sprite) Stop() {
	if s.Cycle == nil {
		return
	} else {
		s.Cycle.Stop()
	}
}