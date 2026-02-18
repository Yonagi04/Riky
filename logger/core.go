package logger

import (
	"io"
	"time"
)

type Core struct {
	levelEnabler LevelEnabler
	encoder      Encoder
	out          io.Writer
}

func NewCore(out io.Writer, encoder Encoder, level Level) *Core {
	return &Core{
		out:          out,
		encoder:      encoder,
		levelEnabler: level,
	}
}

// NewCoreWithEnabler 使用 LevelEnabler 创建 Core
func NewCoreWithEnabler(out io.Writer, encoder Encoder, enabler LevelEnabler) *Core {
	return &Core{
		out:          out,
		encoder:      encoder,
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
	defer putBuffer(buf)

	c.encoder.Encode(buf, msg, level, now, fields)
	_, err := c.out.Write(buf.bs)
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
