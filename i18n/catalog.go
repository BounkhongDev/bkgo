package i18n

// defaultCatalog holds built-in translations for the standard error codes
// defined in the errs package. Add more locales or codes via Register().
var defaultCatalog = map[Locale]map[string]string{
	EN: {
		"NOT_FOUND":      "Resource not found",
		"UNAUTHORIZED":   "Unauthorized",
		"FORBIDDEN":      "Access forbidden",
		"BAD_REQUEST":    "Bad request",
		"CONFLICT":       "Resource already exists",
		"INTERNAL_ERROR": "Internal server error",
		"UNPROCESSABLE":  "Unprocessable entity",
	},
	LO: {
		"NOT_FOUND":      "ບໍ່ພົບຂໍ້ມູນ",
		"UNAUTHORIZED":   "ບໍ່ມີສິດເຂົ້າໃຊ້ງານ",
		"FORBIDDEN":      "ຖືກປະຕິເສດການເຂົ້າໃຊ້ງານ",
		"BAD_REQUEST":    "ຂໍ້ມູນບໍ່ຖືກຕ້ອງ",
		"CONFLICT":       "ຂໍ້ມູນນີ້ມີຢູ່ໃນລະບົບແລ້ວ",
		"INTERNAL_ERROR": "ເກີດຂໍ້ຜິດພາດພາຍໃນລະບົບ",
		"UNPROCESSABLE":  "ບໍ່ສາມາດດຳເນີນການໄດ້",
	},
	TH: {
		"NOT_FOUND":      "ไม่พบข้อมูล",
		"UNAUTHORIZED":   "ไม่มีสิทธิ์เข้าใช้งาน",
		"FORBIDDEN":      "ถูกปฏิเสธการเข้าใช้งาน",
		"BAD_REQUEST":    "ข้อมูลไม่ถูกต้อง",
		"CONFLICT":       "ข้อมูลนี้มีอยู่ในระบบแล้ว",
		"INTERNAL_ERROR": "เกิดข้อผิดพลาดภายในระบบ",
		"UNPROCESSABLE":  "ไม่สามารถดำเนินการได้",
	},
	ZH: {
		"NOT_FOUND":      "未找到资源",
		"UNAUTHORIZED":   "未授权",
		"FORBIDDEN":      "禁止访问",
		"BAD_REQUEST":    "请求错误",
		"CONFLICT":       "资源已存在",
		"INTERNAL_ERROR": "服务器内部错误",
		"UNPROCESSABLE":  "无法处理的实体",
	},
}
