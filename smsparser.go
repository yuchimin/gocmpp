package cmpp
import (
	"fmt"
	"bytes"
	"github.com/bigwhite/gocmpp/utils"
)
/*
	短消息解析器
*/
type SmsParser struct {
	key string	//关键字
	fmt uint8		//编码格式
	total uint8	//分片总数
	len int 		//总字节长度
	msgIds []uint64	//个分片的MsgId
	segments [][]byte //分片数据
}

func (p *SmsParser) reset(key string, fmt uint8, total uint8) {
	p.key = key
	p.fmt = fmt
	p.total = total
	p.len = 0
	p.msgIds = make([]uint64, total)
	p.segments = make([][]byte, total)
}

func (p *SmsParser) Parse(biz string, phoneNum string, msgId uint64, tpUdhi uint8, msgFmt uint8, msgContent string) (string, []uint64) {
	if(tpUdhi==1) {
		//长短信
		buf := []byte(msgContent)
		h_len := uint8(buf[0])
		h_total := uint8(buf[h_len-1])
		h_index := uint8(buf[h_len])

		key := fmt.Sprintf("%s:%s:%d", biz, phoneNum, h_total) //（业务ID:首个手机号:total）作为长短信分片是否属于同一条短信的关键字
		if(p.key!=key) {
			p.reset(key, msgFmt, h_total)
		}
		p.segments[h_index-1] = buf[h_len+1:]
		p.msgIds[h_index-1] = msgId
		p.len += len(p.segments[h_index-1])

		for i := uint8(0); i < h_total; i++ {
			if p.segments[i] == nil {
				return "", nil
			} 
		}
		buf = bytes.Join(p.segments, []byte(""))
		cont, _ := getMsgContent(string(buf), p.fmt)
		p.reset("", 8, 1)
		return cont, p.msgIds
	} else {
		//普通短信
		cont, _ := getMsgContent(msgContent, msgFmt)
		return cont, []uint64{msgId}
	}
}

func getMsgContent(content string, fmt uint8) (string, error) {
	switch fmt {
	case 8:
		return cmpputils.Ucs2ToUtf8(content)
	case 15:
		return cmpputils.GB18030ToUtf8(content)
	default:
		return content, nil
	}
}