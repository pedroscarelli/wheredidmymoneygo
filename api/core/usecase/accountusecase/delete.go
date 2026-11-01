package accountusecase

func (usecase usecase) Delete(accountID string) error {
	return usecase.repository.Delete(accountID)
}
