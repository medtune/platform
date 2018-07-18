package debug

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/medtune/beta-platform/pkg/tmpl"
	"github.com/medtune/beta-platform/pkg/tmpl/data"
)

func Home(c *gin.Context) {
	inject := data.Main{
		PageTitle: "Home",
	}
	c.Status(200)
	tmpl.Home.Execute(c.Writer, &inject)
}

func Index(c *gin.Context) {
	c.Status(200)
	inject := data.Main{
		PageTitle: "Index",
	}
	tmpl.Index.Execute(c.Writer, &inject)
}

func Login(c *gin.Context) {
	c.Status(200)
	tmpl.Login.Execute(c.Writer, &data.Main{
		PageTitle: "login",
	})
}

func Signup(c *gin.Context) {
	c.Status(200)
	tmpl.Signup.Execute(c.Writer, &data.Main{
		PageTitle: "Signup",
	})
}

func Error(c *gin.Context) {
	errStatus := c.Param("code")
	code, err := strconv.Atoi(errStatus)
	if err != nil {
		c.Status(418)
		tmpl.Error.Execute(c.Writer, &data.ErrorImAteaPot)
	} else if code == 404 {
		c.Status(code)
		tmpl.Error.Execute(c.Writer, &data.Error404)
	} else if code == 401 {
		c.Status(401)
		tmpl.Error.Execute(c.Writer, &data.Error401)
	} else if code == 500 {
		c.Status(500)
		tmpl.Error.Execute(c.Writer, &data.Error500)
	} else {
		c.Status(418)
		tmpl.Error.Execute(c.Writer, &data.ErrorFinalBoss)
	}
}

func PolynomialRegression(c *gin.Context) {
	c.Status(200)
	tmpl.DemoPolynomialRegression.Execute(c.Writer, &data.Main{
		PageTitle: "Demo: Polynomial Regression",
	})
}

func ImageClassification(c *gin.Context) {
	c.Status(200)
	tmpl.DemoImageClassification.Execute(c.Writer, &data.Main{
		PageTitle: "Demo: Image classification",
	})
}
