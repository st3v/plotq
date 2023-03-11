package hpgl_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/st3v/plotq/hpgl"
	"github.com/stretchr/testify/require"
)

var svg = `<svg height="50" width="50"><line x1="0" y1="0" x2="50" y2="50" style="stroke:black"/></svg>`

func TestVpypeConvertPortraitSuccess(t *testing.T) {
	expected := "IN;DF;VS10;PS0;SP1;PA;PU0,10870;SP0;IN;\n"

	converter := hpgl.VpypeConverter()

	out, err := converter.Convert(
		bytes.NewReader([]byte(svg)),
		hpgl.Portrait,
		hpgl.Pagesize("a3"),
		hpgl.Device("hp7550"),
		hpgl.Velocity(10),
	)

	require.NoError(t, err)
	require.NotNil(t, out)

	actual, err := ioutil.ReadAll(out)
	require.NoError(t, err)
	require.Equal(t, expected, string(actual))
}

func TestVpypeConvertLandscapeSuccess(t *testing.T) {
	expected := "IN;DF;VS10;PS4;SP1;PA;PU10870,7600;SP0;IN;\n"

	converter := hpgl.VpypeConverter()

	out, err := converter.Convert(
		bytes.NewReader([]byte(svg)),
		hpgl.Landscape,
		hpgl.Pagesize("A4"),
		hpgl.Device("HP7550"),
		hpgl.Velocity(10),
	)

	require.NoError(t, err)
	require.NotNil(t, out)

	actual, err := ioutil.ReadAll(out)
	require.NoError(t, err)
	require.Equal(t, expected, string(actual))
}

func TestVpypeConvertInvalidSVG(t *testing.T) {
	expected := "vpype xml.etree.ElementTree.ParseError: syntax error: line 1, column 0"

	converter := hpgl.VpypeConverter()

	_, err := converter.Convert(
		bytes.NewReader([]byte("invalid")),
		hpgl.Landscape,
		hpgl.Pagesize("a4"),
		hpgl.Device("hp7550"),
		hpgl.Velocity(10),
	)

	require.EqualError(t, err, expected)
}
func TestVpypeConvertInvalidDevice(t *testing.T) {
	expected := "vpype ValueError: no configuration available for plotter 'foo'"

	converter := hpgl.VpypeConverter()

	_, err := converter.Convert(
		bytes.NewReader([]byte(svg)),
		hpgl.Landscape,
		hpgl.Pagesize("a4"),
		hpgl.Device("foo"),
		hpgl.Velocity(10),
	)

	require.EqualError(t, err, expected)
}

func TestVpypeConvertInvalidPagesize(t *testing.T) {
	expected := "vpype ValueError: no configuration available for paper size 'huh' with plotter 'hp7550'"

	converter := hpgl.VpypeConverter()

	_, err := converter.Convert(
		bytes.NewReader([]byte(svg)),
		hpgl.Landscape,
		hpgl.Pagesize("huh"),
		hpgl.Device("hp7550"),
		hpgl.Velocity(10),
	)

	require.EqualError(t, err, expected)
}

func TestVpypeConvertInvalidCommand(t *testing.T) {
	expected := "could not start command: exec: \"invalid\": executable file not found in $PATH"

	converter := hpgl.VpypeConverter(hpgl.VpypeCommand("invalid"))

	_, err := converter.Convert(
		bytes.NewReader([]byte(svg)),
		hpgl.Landscape,
		hpgl.Pagesize("huh"),
		hpgl.Device("hp7550"),
		hpgl.Velocity(10),
	)

	require.EqualError(t, err, expected)
}
