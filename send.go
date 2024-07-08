package ndi

import (
	"log"
	"syscall"
	"unsafe"
)

var (
	ndiLib                 = syscall.NewLazyDLL("Processing.NDI.Lib.x64.dll")
	ndiLibSendCreate       = ndiLib.NewProc("NDIlib_send_create")
	ndiLibSendDestroy      = ndiLib.NewProc("NDIlib_send_destroy")
	ndiLibSendSendVideoV2  = ndiLib.NewProc("NDIlib_send_send_video_v2")
	ndiLibSendGetNoConnections = ndiLib.NewProc("NDIlib_send_get_no_connections")
)

type SendInstance struct{}

func NewSendInstance(settings *SendCreateSettings) *SendInstance {
	log.Println("Creating NDI send instance with settings:", settings)
	ret, _, err := ndiLibSendCreate.Call(uintptr(unsafe.Pointer(settings)))
	if ret == 0 {
		log.Fatalf("Failed to create NDI send instance: %v", err)
	}
	log.Println("NDI send instance created successfully")
	return (*SendInstance)(unsafe.Pointer(ret))
}

func (inst *SendInstance) Destroy() {
	log.Println("Destroying NDI send instance")
	ret, _, err := ndiLibSendDestroy.Call(uintptr(unsafe.Pointer(inst)))
	if ret == 0 {
		if errno, ok := err.(syscall.Errno); ok && errno != 0 {
			log.Fatalf("Failed to destroy NDI send instance: %v", err)
		}
	}
	log.Println("NDI send instance destroyed successfully")
}

func (inst *SendInstance) SendVideoV2(frame *VideoFrameV2) {
	//log.Println("Sending video frame:", frame)
	ret, _, err := ndiLibSendSendVideoV2.Call(uintptr(unsafe.Pointer(inst)), uintptr(unsafe.Pointer(frame)))
	if ret == 0 {
		if errno, ok := err.(syscall.Errno); ok && errno != 0 {
			log.Fatalf("Failed to send video frame: %v", err)
		}
	}
	//log.Println("Video frame sent successfully")
}

func (inst *SendInstance) GetNumConnections(timeoutInMs uint32) (int, error) {
	ret, _, err := ndiLibSendGetNoConnections.Call(uintptr(unsafe.Pointer(inst)), uintptr(timeoutInMs))
	if ret == 0 {
		if errno, ok := err.(syscall.Errno); ok && errno != 0 {
			return 0, Error{errno}
		}
	}
	return int(ret), nil
}
