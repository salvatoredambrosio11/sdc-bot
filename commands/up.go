package commands

import (
	"sdc/models"
	"sdc/service"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var codeBase = make(map[string]models.Codes)

func Up(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	guild, _ := session.Guild(interaction.GuildID)

	if len(interaction.ApplicationCommandData().Options) < 1 {
		captchaKey := service.RandomString(4, "1234567890")
		captchaImage, _ := service.GenerateCaptcha(captchaKey)

		codeBase[interaction.GuildID+interaction.Member.User.ID] = models.Codes{
			Code:      captchaKey,
			Timestamp: time.Now(),
		}

		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Files: []*discordgo.File{
					{
						Name:        "captcha.png",
						ContentType: "image/png",
						Reader:      captchaImage,
					},
				},
				Embeds: []*discordgo.MessageEmbed{
					{
						Description: "Введите число, написанное на изображении, используя команду `/up [код:]`",
						Color:       models.EmbedColors["blue"],
						Author: &discordgo.MessageEmbedAuthor{
							Name:    guild.Name,
							URL:     "https://server-discord.com/" + interaction.GuildID,
							IconURL: guild.IconURL(),
						},
						Footer: &discordgo.MessageEmbedFooter{
							Text: "⏳ Данный код будет действителен в течении 15 секунд!",
						},
						Image: &discordgo.MessageEmbedImage{
							URL: "attachment://captcha.png",
						},
					},
				},
				Flags: 1 << 6,
			},
		})

		return
	}

	enteredCode := interaction.ApplicationCommandData().Options[0].IntValue()

	if generatedCode, ok := codeBase[interaction.GuildID+interaction.Member.User.ID]; ok {

		if time.Now().Sub(generatedCode.Timestamp) > time.Second*15 {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Description: "Срок действия кода истёк!\nПолучите новый, прописав команду `/up`!",
							Color:       models.EmbedColors["red"],
							Author: &discordgo.MessageEmbedAuthor{
								Name:    guild.Name,
								URL:     "https://server-discord.com/" + interaction.GuildID,
								IconURL: guild.IconURL(),
							},
						},
					},
					Flags: 1 << 6,
				},
			})

			return
		}

		if generatedCode.Code == strconv.Itoa(int(enteredCode)) {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Description: "**Успешный Up!**\nВремя фиксации апа: <t:" + strconv.FormatInt(time.Now().Unix(), 10) + ":T>",
							Color:       models.EmbedColors["green"],
							Author: &discordgo.MessageEmbedAuthor{
								Name:    guild.Name,
								URL:     "https://server-discord.com/" + interaction.GuildID,
								IconURL: guild.IconURL(),
							},
						},
					},
				},
			})
			delete(codeBase, interaction.GuildID+interaction.Member.User.ID)
		} else {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Description: "Введённое число не верно!",
							Color:       models.EmbedColors["red"],
							Author: &discordgo.MessageEmbedAuthor{
								Name:    guild.Name,
								URL:     "https://server-discord.com/" + interaction.GuildID,
								IconURL: guild.IconURL(),
							},
						},
					},
					Flags: 1 << 6,
				},
			})
		}
	} else {
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Description: "Введите `/up` без кода, что бы его сгенерировать!",
						Color:       models.EmbedColors["yellow"],
						Author: &discordgo.MessageEmbedAuthor{
							Name:    guild.Name,
							URL:     "https://server-discord.com/" + interaction.GuildID,
							IconURL: guild.IconURL(),
						},
					},
				},
				Flags: 1 << 6,
			},
		})
	}
}
