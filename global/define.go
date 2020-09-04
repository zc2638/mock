/**
 * Created by zc on 2020/9/4.
 */
package global

const DefaultConfigPath = "config.yaml"

var OpenTracingHeaders = []string{
	"x-request-id",
	"x-b3-traceid",
	"x-b3-spanid",
	"x-b3-parentspanid",
	"x-b3-sampled",
	"x-b3-flags",
	"x-ot-span-context",
}
