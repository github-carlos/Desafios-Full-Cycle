package commands

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"
	platformTypes "trevas-bot/pkg/platform/types"
)

type Result struct {
	Labels  []string
	Phrases []string
}

var phrasesWhenIsLow = []string{
	"Apenas afofa o pau do amigo na brotheragem",
	"Só deu o brioco uma vez pra experimentar...",
	"Tem curiosidade mas tem medo de ser julgado",
	"Só tava testando a parada, sem compromisso, saca?",
	"Tava só explorando novas vibes, sem rotular nada.",
	"Experimentou uma vez só pra não morrer curioso.",
	"Curioso, mas não quer que ninguém fique sabendo.",
	"Gosta só de uma resenha entre brothers, nada além.",
	"Essa da uma conferida rápida no documento do brother, só na amizade.",
	"Já deu uma lustrada no boneco do amigo. Apenas uma brincadeira entre amigos, sem significado.",
	"Aprecia um veveco e nada mais",
}

var phrasesWhenIsMedium = []string{
	"Pelo bem e pelo mal, esse cu já levou pau.",
	"Pelo circo e pelo bozo, esse cu já deu gostoso",
	"Pelo som do clarinete, esse cu levou cacete",
	"Pelo som que ele faz, esse cu já deu demais",
	"Pelo ronco do motor, esse cu já fez amor",
	"Pelo som e pelo ronco, seu cu já levou tronco",
	"Pelo toque desse sinto, esse cu levou pepino",
	"Pelo som do clarinete, esse cu levou cacete",
	"Pelo som do piano, seu rabo levou cano",
	"Pelo cheiro da cebola, esse cu já levou rola",
	"Pela barba do profeta, esse cu não tem mais prega",
	"Em chuva de piroca esse aí põe o cu na goteira",
	"Pelo rosnar do cachorro, esse cu pede socorro",
	"Pela Glória Maria, leva pica noite e dia",
	"Pelo cheiro de bacalhau, esse cu agasalha pau",
	"Pelo cheiro da cebola, esse cu já levou rola",
	"Pela som da filosofia, esse cu já deu até cria",
	"Hummmmmmmmmmm boiola",
}

var phrasesWhenIsHigh = []string{
	"Esse da o cu com força!",
	"Esse da o cu que chora.",
	"Em plantação de mandioca esse usa o cu de enxada.",
	"Confirmado!! Esse mama cacete com vontade",
}

var phrasesWhenIsNot = []string{
	"Esse é conhecido como esfolador de buceta",
	"Homem de verdade, odeia viado.",
}

var results = map[string]Result{
	"ultralow": {
		Labels:  []string{"MACHO RAIZ", "HETERO DE VERDADE", "Esfolador de buceta", "Comedor de Cocota", "Rei Delas", "Giga Chad"},
		Phrases: phrasesWhenIsNot,
	},
	"low": {
		Labels:  []string{"Afeminado", "Sojado", "Femboy", "Elu/Delu", "Furry", "Otaku", "Fã de Beyonce", "Escutava Restart"},
		Phrases: phrasesWhenIsLow,
	},
	"medium": {
		Labels:  []string{"Jogador de LoL", "Agasalha Croquete", "Morde fronha", "Toco na Bunda", "Caga Sangue", "Manja Rola", "Belieber", "NPC TikTok", "Faz o pai chorar no banho"},
		Phrases: phrasesWhenIsMedium,
	},
	"high": {Labels: []string{"Ultra Gay", "Viadão", "Boiolão", "Chupingole", "Mamador Profissional", "Fã do Pablo Vittar"}, Phrases: phrasesWhenIsHigh},
}

type ViadometroCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p ViadometroCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Viadometro Command")

  userNumber := extractUserNumber(input.Payload)
  
  if userNumber == "" {
    p.platform.SendReply("Marque alguém para medir o nível de viadagem", &input.EventMessage)
    return
  }

	p.platform.SendReply("⌛⌛ ...Calculando... ⌛⌛", &input.EventMessage)
	time.Sleep(5 * time.Second)

	value := sort(100)

  level := getLevel(value)
  result := results[level]

  randomLabel := getRandomLabel(result)
  randomPhrase := getRandomPhrase(result)
  progressBar := getProgressBar(value)


  fmt.Println("Radom Label", randomLabel)
  fmt.Println("Radom Phrase", randomPhrase)
  fmt.Println("Progressbar", progressBar)
  fmt.Println("User number", userNumber)
  percent := string('%')
  text := fmt.Sprintf("⚠️ Resultado Viadômetro do @%s ⚠️\n\n *0%s %s 100%s*\n Nível: %s\n\n%s", userNumber, percent, progressBar, percent, randomLabel, randomPhrase)

  p.platform.SendText(platformTypes.SendTextInput{Text: text, Mentions: []string{userNumber}}, &input.EventMessage)
}

func (c ViadometroCommand) GetKey() string {
	return c.key
}

func NewViadometroCommand() *ViadometroCommand {
	return &ViadometroCommand{key: "viadometro"}
}

func sort(max int) int {
	return rand.Intn(max)
}

func getLevel(value int) string {
  if value <= 10 {
    return "ultralow"
  }

  if value <= 30 {
    return "low"
  }

  if value <= 70 {
    return "medium"
  }

  return "high"
}

func getProgressBar(value int) string {
  numOfBlocks := value / 10;
  if value <= 10 {
    numOfBlocks = 0;
  }

  blocks := "["

  for i := 0; i <= numOfBlocks; i++ {
    blocks += "■"
  }

  for i := 0; i <  (10 - numOfBlocks); i++ {
    blocks += "•"
  }

  blocks += "]"

  return blocks
}

func getRandomLabel(result Result) string {
	randomInt := rand.Intn(len(result.Labels))
	return result.Labels[randomInt]
}

func getRandomPhrase(result Result) string {
	randomInt := rand.Intn(len(result.Phrases))
	return result.Phrases[randomInt]
}

func extractUserNumber(text string) string {
	regex := regexp.MustCompile(`@(\d+)`)
	result := regex.FindStringSubmatch(text)

	if len(result) > 1 {
		return result[1]
  }
  return ""
}
