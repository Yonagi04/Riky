package logger

import (
	"io"
	"time"

	"github.com/Yonagi04/Riky/encoder"
)

type Core struct {
	levelEnabler LevelEnabler
	encoder      encoder.Encoder
	out          io.Writer
}

func NewCore(out io.Writer, enc encoder.Encoder, level Level) *Core {
	return &Core{
		out:          out,
		encoder:      enc,
		levelEnabler: level,
	}
}

// NewCoreWithEnabler 使用 LevelEnabler 创建 Core
func NewCoreWithEnabler(out io.Writer, enc encoder.Encoder, enabler LevelEnabler) *Core {
	return &Core{
		out:          out,
		encoder:      enc,
		levelEnabler: enabler,
	}
}

func (c *Core) Enabled(level Level) bool {
	return c.levelEnabler.Enabled(level)
}

func (c *Core) Write(level Level, msg string, fields []Field) error {
	if !c.Enabled(level) {
		return nil
	}

	now := time.Now()
	buf := getBuffer()
	c.encoder.Encode(buf, msg, level.String(), now, fields)
	_, err := c.out.Write(buf.Bytes())

	putBuffer(buf)
	return err
}

func (c *Core) With(fields []Field) *Core {
	clone := c.encoder.Clone()
	clone.AddFields(fields)
	return &Core{
		out:          c.out,
		encoder:      clone,
		levelEnabler: c.levelEnabler,
	}
}
