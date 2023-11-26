package Shader

const RendererVertexShader = `
#version 460


in vec2 Vert;
in vec2 Uv;

uniform mat4 Camera;
uniform mat4 Model;

out vec2 a_uv;

void main(){
	a_uv = Uv;
	gl_Position = Camera * Model * vec4(Vert.x,Vert.y,0,1);
}

`

const RendererFragmentShader = `
#version 460

uniform vec4 BaseColor;
uniform sampler2D tex;
uniform bool texEnabled;

in vec2 a_uv;
out vec4 OutputColor;
void main(){
	if(texEnabled){
		OutputColor = BaseColor * texture(tex, a_uv);
	} else{
		OutputColor = BaseColor;
	}
	
}


`
