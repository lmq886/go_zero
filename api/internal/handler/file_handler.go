package handler

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"go_zero/api/internal/model"
	"go_zero/api/internal/svc"
	"go_zero/api/internal/types"
)

type UploadFileHandler struct {
	svcCtx *svc.ServiceContext
}

func NewUploadFileHandler(svcCtx *svc.ServiceContext) *UploadFileHandler {
	return &UploadFileHandler{svcCtx: svcCtx}
}

func (h *UploadFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 获取用户ID
	userId := GetUserIdFromCtx(r.Context())

	// 解析文件
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		logx.Error("Failed to parse multipart form:", err)
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "文件解析失败",
		})
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		logx.Error("Failed to get file:", err)
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "获取文件失败",
		})
		return
	}
	defer file.Close()

	// 检查文件大小
	if handler.Size > h.svcCtx.Config.Upload.MaxSize {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "文件大小超过限制",
		})
		return
	}

	// 检查文件类型
	ext := filepath.Ext(handler.Filename)
	ext = strings.ToLower(ext[1:])
	allowTypes := strings.Split(h.svcCtx.Config.Upload.AllowTypes, ",")
	isAllowed := false
	for _, t := range allowTypes {
		if t == ext {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "不允许的文件类型",
		})
		return
	}

	// 计算MD5
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logx.Error("Failed to read file:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "读取文件失败",
		})
		return
	}

	md5Hash := md5.Sum(fileBytes)
	md5Str := hex.EncodeToString(md5Hash[:])

	// 检查是否已存在相同MD5的文件
	existingFile, _ := h.svcCtx.FileModel.FindByMd5(r.Context(), md5Str)
	if existingFile != nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    0,
			"message": "上传成功",
			"data": map[string]interface{}{
				"id":            existingFile.Id,
				"name":          existingFile.Name,
				"original_name": existingFile.OriginalName,
				"path":          existingFile.Path,
				"url":           existingFile.Url.String,
				"size":          existingFile.Size,
				"type":          existingFile.Type,
				"extension":     existingFile.Extension.String,
				"md5":           existingFile.Md5.String,
				"created_at":    existingFile.CreatedAt.Format("2006-01-02 15:04:05"),
			},
		})
		return
	}

	// 生成存储路径
	now := time.Now()
	datePath := now.Format("2006/01/02")
	savePath := filepath.Join(h.svcCtx.Config.Upload.SavePath, datePath)
	err = os.MkdirAll(savePath, os.ModePerm)
	if err != nil {
		logx.Error("Failed to create directory:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建目录失败",
		})
		return
	}

	// 生成文件名
	fileName := uuid.New().String() + "." + ext
	filePath := filepath.Join(savePath, fileName)

	// 保存文件
	err = os.WriteFile(filePath, fileBytes, 0644)
	if err != nil {
		logx.Error("Failed to save file:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "保存文件失败",
		})
		return
	}

	// 获取文件MIME类型
	contentType := http.DetectContentType(fileBytes)

	// 保存到数据库
	fileModel := &model.File{
		Name:         fileName,
		OriginalName: handler.Filename,
		Path:         filePath,
		Url:          sql.NullString{String: "/uploads/" + datePath + "/" + fileName, Valid: true},
		Size:         handler.Size,
		Type:         contentType,
		Extension:    sql.NullString{String: ext, Valid: true},
		Md5:          sql.NullString{String: md5Str, Valid: true},
		UserId:       sql.NullInt64{Int64: userId, Valid: userId > 0},
	}

	result, err := h.svcCtx.FileModel.Insert(r.Context(), fileModel)
	if err != nil {
		logx.Error("Failed to insert file:", err)
		os.Remove(filePath)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "保存文件信息失败",
		})
		return
	}

	fileId, _ := result.LastInsertId()

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "上传成功",
		"data": map[string]interface{}{
			"id":            fileId,
			"name":          fileName,
			"original_name": handler.Filename,
			"path":          filePath,
			"url":           "/uploads/" + datePath + "/" + fileName,
			"size":          handler.Size,
			"type":          contentType,
			"extension":     ext,
			"md5":           md5Str,
			"created_at":    time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

type GetFileHandler struct {
	svcCtx *svc.ServiceContext
}

func NewGetFileHandler(svcCtx *svc.ServiceContext) *GetFileHandler {
	return &GetFileHandler{svcCtx: svcCtx}
}

func (h *GetFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的文件ID",
		})
		return
	}

	file, err := h.svcCtx.FileModel.FindOne(r.Context(), id)
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "文件不存在",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"id":            file.Id,
			"name":          file.Name,
			"original_name": file.OriginalName,
			"path":          file.Path,
			"url":           file.Url.String,
			"size":          file.Size,
			"type":          file.Type,
			"extension":     file.Extension.String,
			"md5":           file.Md5.String,
			"created_at":    file.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

type ListFilesHandler struct {
	svcCtx *svc.ServiceContext
}

func NewListFilesHandler(svcCtx *svc.ServiceContext) *ListFilesHandler {
	return &ListFilesHandler{svcCtx: svcCtx}
}

func (h *ListFilesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req types.FileListRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err)
		return
	}

	page := int64(req.Page)
	pageSize := int64(req.PageSize)

	files, total, err := h.svcCtx.FileModel.FindPage(r.Context(), page, pageSize, req.Name, req.Type)
	if err != nil {
		logx.Error("Failed to list files:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取文件列表失败",
		})
		return
	}

	list := make([]map[string]interface{}, 0)
	for _, file := range files {
		list = append(list, map[string]interface{}{
			"id":            file.Id,
			"name":          file.Name,
			"original_name": file.OriginalName,
			"path":          file.Path,
			"url":           file.Url.String,
			"size":          file.Size,
			"type":          file.Type,
			"extension":     file.Extension.String,
			"md5":           file.Md5.String,
			"created_at":    file.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"list":      list,
		},
	})
}

type DeleteFileHandler struct {
	svcCtx *svc.ServiceContext
}

func NewDeleteFileHandler(svcCtx *svc.ServiceContext) *DeleteFileHandler {
	return &DeleteFileHandler{svcCtx: svcCtx}
}

func (h *DeleteFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的文件ID",
		})
		return
	}

	file, err := h.svcCtx.FileModel.FindOne(r.Context(), id)
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
			"code":    400,
			"message": "文件不存在",
		})
		return
	}

	// 删除物理文件
	err = os.Remove(file.Path)
	if err != nil {
		logx.Error("Failed to remove file:", err)
	}

	// 删除数据库记录
	err = h.svcCtx.FileModel.Delete(r.Context(), id)
	if err != nil {
		logx.Error("Failed to delete file:", err)
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除文件失败",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "删除成功",
	})
}
