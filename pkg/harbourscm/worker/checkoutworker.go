package worker

type CheckoutWorker struct {
	Github chan GithubCheckoutTask
}

type CheckoutCompletedModel struct {
	WorkspacePath string `json:"string"`
	Success       bool   `json:"success"`
}

func (w CheckoutWorker) DoWork() {
	select {
	case task := <-w.Github:
		CheckoutGithub(task)
	}
}
