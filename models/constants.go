package models

//import "github.com/astaxie/beego"

const (
	C_InitUser                 = 10
	C_PreparePhase             = 0
	C_AllGetReady              = 1
	C_PrepareToStart           = 11
	C_GameStart                = 2
	C_AssignRole               = 21
	C_AssignRoleReady          = 22
	C_LightOffAndShowRole      = 23
	C_LightOffAndShowRoleReady = 24
	C_RoundEnd                 = 33
	C_Round1Start              = 3
	C_Round1End                = 31
	C_Round2Start              = 4
	C_Round2End                = 41
	C_Round3Start              = 5
	C_Round3End                = 51
	C_Round4Start              = 6
	C_Round4End                = 61
	C_Round5Start              = 7
	C_Round5End                = 71
	C_SelectMissioner          = 8
	C_SelectMissionerEnd       = 81
	C_SpeechStart              = 82
	C_SpeechEnd                = 83
	C_VoteStart                = 84
	C_VoteSucess               = 85
	C_VoteFail                 = 851
	C_MissionStart             = 86
	C_MissionEnd               = 87
	C_AssassinateStart         = 88
	C_AssassinateEnd           = 89
	C_ShowGameResult           = 9

	C_UnReady                    = 101
	C_Ready                      = 102
	C_LeaderCanStart             = 103
	C_LeaderCanNotStart          = 1031 //当优信用户加入时通知leader隐藏start键
	C_LeaderStart                = 104
	C_ReStart                    = 1041
	C_GoIntoLightOffPhase        = 105
	C_RolesToMerlin              = 106 //向梅林展示坏人角色
	C_RolesBadToBad              = 107 //坏人互相确认
	C_RolesToPercival            = 108 //向派西维尔展示角色
	C_CountDown                  = 109 //让客户端倒计时
	C_LightOnandR1start          = 110 //天亮，开始第一回合
	C_R2start                    = 111
	C_R3start                    = 112
	C_R4Start                    = 113
	C_R5start                    = 114
	C_LeaderstartSelect          = 115 //leader开始选人
	C_LeaderFinSelect            = 116 //leader选人完毕
	C_StarttoSpeech              = 117 //
	C_StartToVote                = 118
	C_VoteFin                    = 119
	C_VoteFailAndReSelect        = 120
	C_VoteSuccAndStartMission    = 121
	C_DelayOver4TimesAndNoSpeech = 122
	C_StartToMission             = 123
	C_MissionResult              = 124
	C_Assassinate                = 125
	C_AssassinateResult          = 126

	CHAT_MTYPE    = "chat_mtype"
	NOTIFY_MTYPE  = "notify_mtype"
	CONTROL_MTYPE = "control_mtype"
	UPDATE_MTYPE  = "update_mtype"

	TIME_FORMAT = "01-02 15:04:05"

	//role
	//梅林、派西维尔、忠臣*4  vs 莫德雷德、莫甘娜、奥伯伦、刺客
	R_MERLIN   = 1
	R_Percival = 2
	R_loyalist = 3
	R_Mordred  = 4
	R_Morgana  = 5
	R_Oberon   = 6
	R_Assassin = 7
	R_Minion   = 8
	R_Unknown  = 9
)
