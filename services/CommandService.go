package services

import (
  "crypto/md5"
  "deploy-server/model"
  "deploy-server/utils"
  "encoding/hex"
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "os"
)

func RunWrapperCommand(command string) model.CommandResult {
  result := model.CommandResult{}
  if "nginx-reload" == command {
    result = utils.RunComamnd("nginx", "-s", "reload")
  } else if "nginx-t" == command {
    result = utils.RunComamnd("nginx", "-t")
  } else {
    //创建一个文件
    hash := md5.New()
    hash.Write([]byte(command))
    md5String := hex.EncodeToString(hash.Sum(nil))
    movedDir := "script"
    // 创建移动路径
    err := os.MkdirAll(movedDir, 0755)
    if err != nil {
      result.Error = err.Error()
      result.Success = false
      return result

    }
    //移动到这个文件夹
    dstFilePath := movedDir + "/" + md5String
    hlog.Info("dstFilePath:", dstFilePath)
    //创建文件
    scriptFile, err := os.OpenFile(dstFilePath, os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
      result.Error = err.Error()
      result.Success = false
      return result
    }
    defer scriptFile.Close()
    n, err := scriptFile.WriteString(command)
    if err != nil {
      result.Error = err.Error()
      result.Success = false
      return result
    } else {
      hlog.Info(n)
    }
    hlog.Info("run script:", "sh", dstFilePath)
    result = utils.RunComamnd("sh", dstFilePath)
  }
  return result
}
