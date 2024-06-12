package utils

import (
	"bc-alert/src/models"
	"fmt"
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func BuildFlexContainer(todayPrice models.GoldPriceData) (*linebot.FlexMessage, error) {
	smSize := linebot.FlexTextSizeTypeXxs
	xxlSize := linebot.FlexTextSizeTypeXxl
	mdMargin := linebot.FlexComponentMarginTypeMd
	xxlMargin := linebot.FlexComponentMarginTypeXxl

	barBuy := fmt.Sprintf("%v", todayPrice.BarBuy)
	barSell := fmt.Sprintf("%v", todayPrice.BarSell)
	ornamentBuy := fmt.Sprintf("%v", todayPrice.OrnamentBuy)
	ornamentSell := fmt.Sprintf("%v", todayPrice.OrnamentSell)
	dateDay := fmt.Sprintf(" %v\t%v\t%v", todayPrice.UpdatedDate, todayPrice.UpdatedTime, todayPrice.UpdateTheTime)

	colorCode := "#666666"
	emoji := ""
	statusChange, err := strconv.Atoi(todayPrice.TodayChange)
	if err != nil {
	} else {
		if statusChange > 0 {
			colorCode = "#27ae60"
			emoji = "⬆"
		} else if statusChange < 0 {
			colorCode = "#c0392b"
			emoji = "⬇"
		}
	}

	toDay := fmt.Sprintf("วันนี้ %v %v", emoji, todayPrice.TodayChange)

	bubble := linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeGiga,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "ราคาทองตามประกาศสมาคมค้าทองคำ",
					Size:   smSize,
					Color:  "#aaaaaa",
					Weight: linebot.FlexTextWeightTypeBold,
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   toDay,
					Size:   xxlSize,
					Color:  colorCode,
					Weight: linebot.FlexTextWeightTypeBold,
					Margin: mdMargin,
				},
				&linebot.TextComponent{
					Type:  linebot.FlexComponentTypeText,
					Text:  dateDay,
					Size:  linebot.FlexTextSizeTypeXs,
					Color: "#aaaaaa",
					Wrap:  true,
				},
				&linebot.SeparatorComponent{
					Type:   linebot.FlexComponentTypeSeparator,
					Margin: xxlMargin,
				},
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeVertical,
					Margin:  xxlMargin,
					Spacing: linebot.FlexComponentSpacingTypeSm,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeHorizontal,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "ทองแท่ง-รับซื้อ",
									Size:  smSize,
									Color: "#555555",
									Flex:  &[]int{0}[0],
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  barBuy,
									Size:  smSize,
									Color: "#111111",
									Align: linebot.FlexComponentAlignTypeEnd,
								},
							},
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeHorizontal,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "ทองแท่ง-ขายออก",
									Size:  smSize,
									Color: "#555555",
									Flex:  &[]int{0}[0],
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  barSell,
									Size:  smSize,
									Color: "#111111",
									Align: linebot.FlexComponentAlignTypeEnd,
								},
							},
						},
						&linebot.SeparatorComponent{
							Type:   linebot.FlexComponentTypeSeparator,
							Margin: xxlMargin,
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeHorizontal,
							Margin: xxlMargin,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "ทองรูปพรรณ-รับซื้อ",
									Size:  smSize,
									Color: "#555555",
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  ornamentBuy,
									Size:  smSize,
									Color: "#111111",
									Align: linebot.FlexComponentAlignTypeEnd,
								},
							},
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeHorizontal,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "ทองรูปพรรณ-ขายออก",
									Size:  smSize,
									Color: "#555555",
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  ornamentSell,
									Size:  smSize,
									Color: "#111111",
									Align: linebot.FlexComponentAlignTypeEnd,
								},
							},
						},
					},
				},
			},
		},
		Styles: &linebot.BubbleStyle{
			Footer: &linebot.BlockStyle{
				Separator: true,
			},
		},
	}

	flexMessage := linebot.NewFlexMessage("Update ราคาทองวันนี้", &bubble)

	return flexMessage, nil
}
