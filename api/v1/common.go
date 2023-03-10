package v1

const (
	DefaultVelocity    = 50
	DefaultOrientation = OrientationPortrait
)

type Orientation string

const (
	OrientationLandscape Orientation = "landscape"
	OrientationPortrait  Orientation = "portrait"
)

func (Orientation) Enum() []interface{} {
	return []interface{}{
		OrientationLandscape,
		OrientationPortrait,
	}
}

type Pagesize string

const (
	PagesizeA0        Pagesize = "a0"
	PagesizeA1        Pagesize = "a1"
	PagesizeA2        Pagesize = "a2"
	PagesizeA3        Pagesize = "a3"
	PagesizeA4        Pagesize = "a4"
	PagesizeA5        Pagesize = "a5"
	PagesizeA6        Pagesize = "a6"
	PagesizeExecutive Pagesize = "executive"
	PagesizeLegal     Pagesize = "legal"
	PagesizeLetter    Pagesize = "letter"
	PagesizeTabloid   Pagesize = "tabloid"
	PagesizeTight     Pagesize = "tight"
)

func (Pagesize) Enum() []interface{} {
	return []interface{}{
		PagesizeA0,
		PagesizeA1,
		PagesizeA2,
		PagesizeA3,
		PagesizeA4,
		PagesizeA5,
		PagesizeA6,
		PagesizeExecutive,
		PagesizeLegal,
		PagesizeLetter,
		PagesizeTabloid,
		PagesizeTight,
	}
}

type Device string

const (
	DeviceArtisan    Device = "artisan"
	DeviceDesignmate Device = "designmate"
	DeviceDMP161     Device = "dmp_161"
	DeviceDXY        Device = "dxy"
	DeviceHP7475A    Device = "hp7475a"
	DeviceHP7440A    Device = "hp7440a"
	DeviceHP7550     Device = "hp7550"
	DeviceSketchmate Device = "sketchmate"
)

func (Device) Enum() []interface{} {
	return []interface{}{
		DeviceArtisan,
		DeviceDesignmate,
		DeviceDMP161,
		DeviceDXY,
		DeviceHP7475A,
		DeviceHP7440A,
		DeviceHP7550,
		DeviceSketchmate,
	}
}
