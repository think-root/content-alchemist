package server

import "testing"

func TestValidateGeneratedDescription(t *testing.T) {
	const validUk = "Проєкт є локальним проксі, що зменшує витрати на LLM, рендерячи громіздкий контекст у компактні зображення та використовуючи можливості vision для стиснення системних промптів і великої документації в один запит користувача."

	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{name: "empty", text: "", wantErr: true},
		{name: "refusal", text: "Будь ласка, надайте вміст файлу README.md, оскільки посилання на відео не дозволяє мені проаналізувати технічну документацію репозиторію.", wantErr: true},
		{name: "too short", text: "Короткий опис проєкту.", wantErr: true},
		{name: "valid single", text: validUk, wantErr: false},
		{name: "valid multilingual", text: "===(uk)" + validUk + "===", wantErr: false},
		{name: "multilingual with refusal segment", text: "===(uk)" + validUk + "===(en)I cannot analyze this repository without a README.===", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGeneratedDescription(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGeneratedDescription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
