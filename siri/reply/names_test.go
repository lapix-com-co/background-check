package reply

import "testing"

func TestNames_Answer(t *testing.T) {
	type args struct {
		in0 User
		in1 string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   bool
		wantErr bool
	}{
		{
			name:    "should return the users first name",
			args:    args{
				in0: User{
					Name: Name{FirstName: "John"},
				},
				in1: "¿Cual es el primer nombre de la persona a la cual esta expidiendo el certificado?",
			},
			want:    "John",
			want1:   true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Names{}
			got, got1, err := n.Answer(tt.args.in0, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Answer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Answer() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Answer() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNames_Is(t *testing.T) {
	type args struct {
		in0 User
		in1 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should mark the question as valid",
			args: args{
				in0: User{},
				in1: "¿Cual es el primer nombre de la persona a la cual esta expidiendo el certificado?",
			},
			want: true,
		},
		{
			name: "should mark the question as not valid",
			args: args{
				in0: User{},
				in1: "¿Cual son las primeras tres letras del primer nombre de la persona a la cual esta expidiendo el certificado?",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Names{}
			if got := n.Is(tt.args.in0, tt.args.in1); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}