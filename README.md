你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# *fpd.Decimal [![Build Status](https://travis-ci.org/oguzbilgic/fpd.png?branch=master)](https://travis-ci.org/oguzbilgic/fpd)

Package implements fixed-point decimal 

## Usage

```go
package main

import "github.com/oguzbilgic/fpd"

func main() {
	// Buy price of the security: $136.02
	buyPrice := fpd.New(13602000, -5)

	// Sell price of the security: $137.699
	sellPrice := fpd.New(13769900, -5)

	// Volume traded: 0.01
	volume := fpd.New(1000000, -8)

	// Trade fee percentage: 0.6%
	feePer := fpd.New(6, -3)

	buyCost := buyPrice.Mul(volume)
	buyFee := buyPrice.Mul(volume).Mul(feePer)
	sellRevenue := sellPrice.Mul(volume)
	sellFee := sellPrice.Mul(volume).Mul(feePer)

	// Initall account balance: $2.00000
	balance := fpd.New(200000, -5)

	balance = balance.Sub(buyCost)
	balance = balance.Sub(buyFee)
	balance = balance.Add(sellRevenue)
	balance = balance.Sub(sellFee)

	// Final balance
	fmt.Println(balance)
	// Did this trade turn into profit? :)
}
```

## Documentation

http://godoc.org/github.com/oguzbilgic/fpd

## License

The MIT License (MIT)
