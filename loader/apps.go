package loader

import "funny/attack"
import "funny/baidu_pan"
import "funny/exame_result"
import "funny/food"
import "funny/leet_code"
import "funny/my_charles"
import "funny/proxy"
import "funny/spider_client"
import "funny/tcp_log"
import "funny/tencent_code"
import "funny/test"
import "funny/wasm"
import "funny/xes"
import "funny/yeb_exp"


func init() {
	mp = map[string]initFunc {
        "attack": attack.Init,
        "baidu_pan": baidu_pan.Init,
        "exame_result": exame_result.Init,
        "food": food.Init,
        "leet_code": leet_code.Init,
        "my_charles": my_charles.Init,
        "proxy": proxy.Init,
        "spider_client": spider_client.Init,
        "tcp_log": tcp_log.Init,
        "tencent_code": tencent_code.Init,
        "test": test.Init,
        "wasm": wasm.Init,
        "xes": xes.Init,
        "yeb_exp": yeb_exp.Init,

	}
}
