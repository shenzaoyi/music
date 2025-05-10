package tengcent_cos

import (
	"Music/config"
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
)

type CosClient struct {
	client *cos.Client
}

func InitClient() (*CosClient, error) {
	// 将 examplebucket-1250000000 和 COS_REGION 修改为用户真实的信息
	// 存储桶名称，由 bucketname-appid 组成，appid 必须填入，可以在 COS 控制台查看存储桶名称。https://console.cloud.tencent.com/cos5/bucket
	// COS_REGION 可以在控制台查看，https://console.cloud.tencent.com/cos5/bucket, 关于地域的详情见 https://cloud.tencent.com/document/product/436/6224
	u, _ := url.Parse("https://music-1320864532.cos.ap-guangzhou.myqcloud.com")
	// 用于 Get Service 查询，默认全地域 service.cos.myqcloud.com
	su, _ := url.Parse("https://cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u, ServiceURL: su}
	// 1.永久密钥
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.Config.TencentCOS.SecretID,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
			SecretKey: config.Config.TencentCOS.SecretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
		},
	})
	//s, _, err := client.Service.Get(context.Background())
	//if err != nil {
	//	panic(err)
	//}
	//for _, b := range s.Buckets {
	//	fmt.Printf("%#v\n", b)
	//}
	return &CosClient{
		client: client,
	}, nil
}

func (cosClient *CosClient) Upload(name string, filepath string) string {
	res, _, err := cosClient.client.Object.Upload(context.Background(), name, filepath, nil)
	if err != nil {
		fmt.Println("error")
		panic(err)
	}
	return res.Location
}
func (cosClient *CosClient) DownloadStream(key string) (io.ReadCloser, error) {
	resp, err := cosClient.client.Object.Get(context.Background(), key, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
