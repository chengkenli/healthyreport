/*
 *@author  chengkenli
 *@project SchemaHealthy
 *@package tools
 *@file    emailbody
 *@date    2024/7/26 10:08
 */

package tools

import "fmt"

type EsBody struct {
    App        string
    Zty        string
    Num        int
    TableClass string
}

func Ebody(e *EsBody) string {
    emsg := fmt.Sprintf(`
<p>
	<br />
</p>
<p>
	尊敬的starrocks[主题域] 同事:
</p>
<p>
	<span>StarRocks %s集群Tablet治理，有很多不规范的数据表，需要大家协助治理！</span>
</p>
	<p>治理选项包括了【副本异常】、【分桶异常】、【空存储分区过多】三个。</p>
	<p>其中如何确定分桶数量这是一个难题。但这里每个表如何更改，都会在备注栏给大家最详细精准的说明，请按照指引进行改造！</p>
	<p>如果您对分桶这个有兴趣，您可以参考文档：<a href="https://xxxxxxxxxxxx" target="_blank">BUCKETS</a>，避免后面继续创建不标准的属性表。</p>
<p>
	<span><br />
</span> 
</p>
<p>
	总体报告展示：
</p>
<p>
	<table style="width:80%%;" cellpadding="0" cellspacing="0" border="1" bordercolor="#000000" class="ke-zeroborder">
		<tbody>
			<tr>
				<td>
					集群名称<br />
				</td>
				<td>
					主题域<br />
				</td>
				<td>
					责任人<br />
				</td>
				<td>
					需治理表数量<br />
				</td>
				<td>
					跳转链接(获取报告文件)<br />
				</td>
			</tr>
			%s
		</tbody>
	</table>
</p>
<p>
	<br />
</p>
<p>
	此电子邮件由系统自动发送，请勿直接回复此电子邮件。<span> </span> 
</p>
<p>
	谢谢！<span> </span> 
</p>
<p>
	StarRocks
</p>
<p>
	<br />
</p>
<br />
<br />
`, e.App, e.TableClass)
    return emsg
}
