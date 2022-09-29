

package midlevel

func Do() error{
	err := lowlevel.Execute(cmd)
	if err != nil {
		// long error msg string, need to consistent in every level, no call stack
	    return fmt.Errorf("get a error in midlevel Do: %s", err)
	    // sometimes we do: 
	    // lose context of root cause and error flow
	    return fmt.Errorf("get a error in midlevel Do")
	    // or
	    // lose context of midlevel
	    return err
	}
}
