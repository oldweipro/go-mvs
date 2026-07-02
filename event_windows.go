//go:build windows && amd64

package mvs

import (
	"fmt"
	"syscall"
	"unsafe"
)

func (c *Camera) RegisterAllEventCallback(callback EventCallback) error {
	if callback == nil {
		return fmt.Errorf("%w: event callback is nil", ErrInvalidArgument)
	}

	callbackPtr := newEventCallback(callback)

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	if err := c.sdk.driver.registerAllEventCallback(c.handle, callbackPtr, 0); err != nil {
		return err
	}
	c.eventPtrs = append(c.eventPtrs, callbackPtr)
	c.eventCallbacks = append(c.eventCallbacks, callback)
	return nil
}

func (c *Camera) RegisterEventCallback(eventName string, callback EventCallback) error {
	if eventName == "" {
		return fmt.Errorf("%w: event name is empty", ErrInvalidArgument)
	}
	if callback == nil {
		return fmt.Errorf("%w: event callback is nil", ErrInvalidArgument)
	}

	callbackPtr := newEventCallback(callback)

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	if err := c.sdk.driver.registerEventCallbackEx(c.handle, eventName, callbackPtr, 0); err != nil {
		return err
	}
	c.eventPtrs = append(c.eventPtrs, callbackPtr)
	c.eventCallbacks = append(c.eventCallbacks, callback)
	return nil
}

func (c *Camera) EventNotificationOn(eventName string) error {
	if eventName == "" {
		return fmt.Errorf("%w: event name is empty", ErrInvalidArgument)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.eventNotificationOn(c.handle, eventName)
}

func (c *Camera) EventNotificationOff(eventName string) error {
	if eventName == "" {
		return fmt.Errorf("%w: event name is empty", ErrInvalidArgument)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.open || c.handle == 0 {
		return ErrCameraClosed
	}
	return c.sdk.driver.eventNotificationOff(c.handle, eventName)
}

func newEventCallback(callback EventCallback) uintptr {
	return syscall.NewCallback(func(raw *mvEventOutInfo, user uintptr) uintptr {
		defer func() {
			_ = recover()
		}()
		if raw == nil {
			return 0
		}
		callback(eventInfoFromRaw(raw))
		return 0
	})
}

func eventInfoFromRaw(raw *mvEventOutInfo) EventInfo {
	info := EventInfo{
		Name:          byteString(raw.EventName[:]),
		EventID:       raw.EventID,
		StreamChannel: raw.StreamChannel,
		BlockID:       uint64(raw.BlockIDHigh)<<32 | uint64(raw.BlockIDLow),
		Timestamp:     uint64(raw.TimestampHigh)<<32 | uint64(raw.TimestampLow),
	}
	if raw.EventData != nil && raw.EventDataSize > 0 {
		maxInt := uint64(int(^uint(0) >> 1))
		if uint64(raw.EventDataSize) > maxInt {
			return info
		}
		data := unsafe.Slice(raw.EventData, int(raw.EventDataSize))
		info.Data = make([]byte, len(data))
		copy(info.Data, data)
	}
	return info
}
