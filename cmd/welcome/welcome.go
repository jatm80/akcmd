package welcome

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akcmd/l10n"
)

func Cmd(app *gcli.App) *gcli.Command {
	localizedStrings := l10n.GetLocalizationStrings()

	cmd := &gcli.Command{
		Name: "welcome",
		// allow color tag and {$cmd} will be replace to 'demo'
		Desc: localizedStrings.Command["welcome"],
		Func: func(cmd *gcli.Command, args []string) error {
			color.Redp(`
          ///////////         
            //////////        
             ///////////      
               //////////     
    /////////// //////////,   
   //////////.    //////////  
 ///////////       ////////// 
 /////////***********/////////
  ///////   **********/////// 
    ///*     ***********///   
     /         **********/` + "\n\n")

			color.Yellowln(localizedStrings.ClientTitle + ` ` + app.Version)
			color.Info.Println(localizedStrings.WelcomeFirstRun + "\n\n")
			return nil
		},
	}

	return cmd
}
