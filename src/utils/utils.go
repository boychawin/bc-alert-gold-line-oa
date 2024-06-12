package utils

import (
	"bc-alert/src/models"
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"
)

func FlexContainerGenerator(todayPrice models.GoldPriceData) (*linebot.FlexMessage, error) {

	flexMessage := linebot.NewFlexMessage("Gold Price Today", &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Hero: &linebot.ImageComponent{
			Type:        linebot.FlexComponentTypeImage,
			URL:         "https://i.postimg.cc/3xxWh7Ln/bullion-1744773-1280.jpg",
			Size:        linebot.FlexImageSizeTypeFull,
			AspectRatio: linebot.FlexImageAspectRatioType20to13,
			AspectMode:  linebot.FlexImageAspectModeTypeCover,
			Action:      linebot.NewURIAction("Action", "https://linecorp.com/"),
		},
		Body: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: linebot.FlexComponentSpacingTypeMd,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:    linebot.FlexComponentTypeText,
					Text:    "ราคาทองวันนี้",
					Size:    linebot.FlexTextSizeTypeXl,
					Gravity: linebot.FlexComponentGravityTypeCenter,
					Weight:  linebot.FlexTextWeightTypeBold,
					Wrap:    true,
				},
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeVertical,
					Spacing: linebot.FlexComponentSpacingTypeSm,
					Margin:  linebot.FlexComponentMarginTypeLg,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeBaseline,
							Spacing: linebot.FlexComponentSpacingTypeSm,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type: linebot.FlexComponentTypeText,
									Text: "ทองคำแท่ง (ซื้อ)",
									// Flex:  4,
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#AAAAAA",
								},
								&linebot.TextComponent{
									Type: linebot.FlexComponentTypeText,
									Text: fmt.Sprintf("%v", todayPrice.BarBuy),
									// Flex:  3,
									Size:  linebot.FlexTextSizeTypeSm,
									Align: linebot.FlexComponentAlignTypeEnd,
									Color: "#666666",
									Wrap:  true,
								},
							},
						},
					},
				},
			},
		},
	})

	return flexMessage, nil
}
