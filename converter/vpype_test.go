package converter_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/st3v/plotq/converter"
	"github.com/stretchr/testify/require"
)

var svg = `<svg height="50" width="50"><line x1="0" y1="0" x2="50" y2="50" style="stroke:black"/></svg>`

func TestVpypeConvertPortraitSuccess(t *testing.T) {
	expected := "IN;DF;VS10;PS0;SP1;PA;PU0,10870;SP0;IN;\n"

	vpype := converter.Vpype()

	w := vpype.Convert(
		bytes.NewReader([]byte(svg)),
		converter.Portrait,
		converter.Pagesize("a3"),
		converter.Device("hp7550"),
		converter.Velocity(10),
	)

	out := &bytes.Buffer{}
	n, err := w.WriteTo(out)
	require.NoError(t, err)
	require.Equal(t, int64(len(expected)), n)

	actual, err := io.ReadAll(out)
	require.NoError(t, err)
	require.Equal(t, expected, string(actual))
}

func TestVpypeConvertLandscapeSuccess(t *testing.T) {
	expected := "IN;DF;VS10;PS4;SP1;PA;PU10870,7600;SP0;IN;\n"

	vpype := converter.Vpype()

	w := vpype.Convert(
		bytes.NewReader([]byte(svg)),
		converter.Landscape,
		converter.Pagesize("A4"),
		converter.Device("HP7550"),
		converter.Velocity(10),
	)

	out := &bytes.Buffer{}
	n, err := w.WriteTo(out)
	require.NoError(t, err)
	require.Equal(t, int64(len(expected)), n)

	actual, err := io.ReadAll(out)
	require.NoError(t, err)
	require.Equal(t, expected, string(actual))
}

func TestVpypeConvertInvalidSVG(t *testing.T) {
	expected := "ParseError"

	vpype := converter.Vpype()

	w := vpype.Convert(
		bytes.NewReader([]byte("invalid")),
		converter.Landscape,
		converter.Pagesize("a4"),
		converter.Device("hp7550"),
		converter.Velocity(10),
	)

	out := &bytes.Buffer{}
	_, err := w.WriteTo(out)
	require.ErrorContains(t, err, expected)
}
func TestVpypeConvertInvalidDevice(t *testing.T) {
	expected := "no configuration available for plotter 'foo'"

	vpype := converter.Vpype()

	w := vpype.Convert(
		bytes.NewReader([]byte(svg)),
		converter.Landscape,
		converter.Pagesize("a4"),
		converter.Device("foo"),
		converter.Velocity(10),
	)

	out := &bytes.Buffer{}
	_, err := w.WriteTo(out)
	require.ErrorContains(t, err, expected)
}

func TestVpypeConvertInvalidPagesize(t *testing.T) {
	expected := "no configuration available for paper size 'huh'"

	vpype := converter.Vpype()

	w := vpype.Convert(
		bytes.NewReader([]byte(svg)),
		converter.Landscape,
		converter.Pagesize("huh"),
		converter.Device("hp7550"),
		converter.Velocity(10),
	)

	out := &bytes.Buffer{}
	_, err := w.WriteTo(out)
	require.ErrorContains(t, err, expected)
}

func TestVpypeConvertInvalidCommand(t *testing.T) {
	expected := "could not start command"

	vpype := converter.Vpype(converter.VpypeCommand("invalid"))

	w := vpype.Convert(
		bytes.NewReader([]byte(svg)),
		converter.Landscape,
		converter.Pagesize("huh"),
		converter.Device("hp7550"),
		converter.Velocity(10),
	)

	out := &bytes.Buffer{}
	_, err := w.WriteTo(out)
	require.ErrorContains(t, err, expected)
}
