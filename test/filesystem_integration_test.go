package test_test

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/code"
	filesystem "github.com/TencentCloudAgentRuntime/ags-go-sdk/tool/filesystem"
)

func TestFilesystem_GetInfo_Root(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		info, err := sb.Files.GetInfo(context.TODO(), "/", nil)
		if err != nil {
			t.Fatalf("GetInfo('/') error: %v", err)
		}
		if info == nil || info.Name == "" {
			t.Fatalf("invalid root info: %+v", info)
		}
	})
}

func TestFilesystem_MakeDir_List_Remove(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		uniq := time.Now().UnixNano()
		dirPath := "/home/user/fs-integ-dir-" + fmtInt64(uniq)

		created, err := sb.Files.MakeDir(context.TODO(), dirPath, &filesystem.MakeDirConfig{User: "user"})
		if err != nil {
			t.Fatalf("MakeDir error: %v", err)
		}
		if !created {
			t.Logf("MakeDir returned false (may be idempotent), path=%s", dirPath)
		}

		exists, err := sb.Files.Exists(context.TODO(), dirPath, nil)
		if err != nil {
			t.Fatalf("Exists error: %v", err)
		}
		if !exists {
			t.Fatalf("dir should exist after MakeDir: %s", dirPath)
		}

		entries, err := sb.Files.List(context.TODO(), "/home/user", &filesystem.ListConfig{Depth: 1})
		if err != nil {
			t.Fatalf("List error: %v", err)
		}
		found := false
		for _, e := range entries {
			if e.Path == dirPath {
				found = true
				break
			}
		}
		if !found {
			t.Logf("created dir not listed in depth=1; entries=%d", len(entries))
		}

		// 清理
		if err := sb.Files.Remove(context.TODO(), dirPath, nil); err != nil {
			t.Fatalf("Remove dir error: %v", err)
		}
		exists, err = sb.Files.Exists(context.TODO(), dirPath, nil)
		if err != nil {
			t.Fatalf("Exists after remove error: %v", err)
		}
		if exists {
			t.Fatalf("dir still exists after remove: %s", dirPath)
		}
	})
}

func TestFilesystem_Write_Read_Rename_Remove(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		uniq := time.Now().UnixNano()
		filePath := "/home/user/fs-integ-file-" + fmtInt64(uniq) + ".txt"
		newPath := "/home/user/fs-integ-file-" + fmtInt64(uniq) + "-renamed.txt"
		content := "hello-filesystem"

		winfo, err := sb.Files.Write(context.TODO(), filePath, strings.NewReader(content), &filesystem.WriteConfig{User: "user"})
		if err != nil {
			t.Fatalf("Write error: %v", err)
		}
		if winfo == nil || winfo.Path == "" {
			t.Fatalf("invalid write info: %+v", winfo)
		}

		r, err := sb.Files.Read(context.TODO(), filePath, nil)
		if err != nil {
			t.Fatalf("Read error: %v", err)
		}
		data, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("ReadAll error: %v", err)
		}
		if string(data) != content {
			t.Fatalf("unexpected file content: got=%q want=%q", string(data), content)
		}

		if err := sb.Files.Rename(context.TODO(), filePath, newPath, nil); err != nil {
			t.Fatalf("Rename error: %v", err)
		}

		existsNew, err := sb.Files.Exists(context.TODO(), newPath, nil)
		if err != nil {
			t.Fatalf("Exists(new) error: %v", err)
		}
		if !existsNew {
			t.Fatalf("renamed file not found: %s", newPath)
		}
		existsOld, err := sb.Files.Exists(context.TODO(), filePath, nil)
		if err != nil {
			t.Fatalf("Exists(old) error: %v", err)
		}
		if existsOld {
			t.Fatalf("old file still exists after rename: %s", filePath)
		}

		if err := sb.Files.Remove(context.TODO(), newPath, nil); err != nil {
			t.Fatalf("Remove renamed file error: %v", err)
		}
	})
}

func TestFilesystem_Exists_Read_Nonexistent(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		path := "/home/user/not-exist-" + fmtInt64(time.Now().UnixNano())
		exists, err := sb.Files.Exists(context.TODO(), path, nil)
		if err != nil {
			t.Fatalf("Exists(nonexistent) error: %v", err)
		}
		if exists {
			t.Fatalf("nonexistent path reported as exists: %s", path)
		}
		_, err = sb.Files.Read(context.TODO(), path, nil)
		if err == nil {
			t.Fatalf("expected read error for nonexistent path")
		}
	})
}
