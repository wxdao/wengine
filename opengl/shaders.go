package opengl

var defaultMeshShader_COLOR_NOLIGHT = glShaderProgram{
	vertexSource: `
		#version 410 core

		layout (location = 0) in vec3 position;
		layout (location = 1) in vec2 uv;
		layout (location = 2) in vec3 normal;

		uniform mat4 model;
		uniform mat4 view;
		uniform mat4 projection;
		uniform vec4 color;

		out vec4 vs_color;

		void main() {
			vs_color = color;
			gl_Position = projection * view * model * vec4(position, 1.0);
		}
	`,
	fragmentSource: `
		#version 410 core

		in vec4 vs_color;

		out vec4 color;

		void main() {
			color = vs_color;
		}
	`,
}

var defaultMeshShader_TEXTURE_NOLIGHT = glShaderProgram{
	vertexSource: `
		#version 410 core

		layout (location = 0) in vec3 position;
		layout (location = 1) in vec2 uv;
		layout (location = 2) in vec3 normal;

		uniform mat4 model;
		uniform mat4 view;
		uniform mat4 projection;

		out vec2 vs_uv;

		void main() {
			vs_uv = uv;
			gl_Position = projection * view * model * vec4(position, 1.0);
		}
	`,
	fragmentSource: `
		#version 410 core

		in vec2 vs_uv;

		uniform sampler2D diffuseMap;

		out vec4 color;

		void main() {
			color = texture(diffuseMap, vs_uv);
		}
	`,
}

// ----------------------------------------------------------------------------------------------

var defaultMeshShader_COLOR = glShaderProgram{
	vertexSource: `
		#version 410 core

		layout (location = 0) in vec3 position;
		layout (location = 1) in vec2 uv;
		layout (location = 2) in vec3 normal;

		uniform mat4 model;
		uniform mat4 view;
		uniform mat4 projection;
		uniform vec4 color;

		out vec4 vs_color;
		out vec3 vs_normal;
		out vec3 vs_fragPosition;

		void main() {
			vs_color = color;
			vs_normal = mat3(transpose(inverse(model))) * normal;
			vs_fragPosition = vec3(model * vec4(position, 1.0));
			gl_Position = projection * view * model * vec4(position, 1.0);
		}
	`,
	fragmentSource: `
		#version 410 core

		struct dirLight {
			vec3 position;
			vec3 direction;
			vec3 ambient;
			vec3 diffuse;
			vec3 specular;
		};

		struct pointLight {
			vec3 position;
			float range;
			vec3 ambient;
			vec3 diffuse;
			vec3 specular;
		};

		in vec4 vs_color;
		in vec3 vs_normal;
		in vec3 vs_fragPosition;

		uniform vec3 ambient;

		uniform dirLight dirLights[10];
		uniform int num_dirLight = 0;

		uniform pointLight pointLights[50];
		uniform int num_pointLight = 0;

		uniform vec3 cameraPosition;

		out vec4 color;

		vec3 calculateDirLight(dirLight light) {
			vec3 meshDiffuse = vec3(vs_color);
			vec3 viewDirection = normalize(vs_fragPosition - cameraPosition);
			vec3 halfDirection = -normalize(viewDirection + light.direction);

			vec3 diffuse = light.diffuse * max(dot(normalize(vs_normal), normalize(-light.direction)), 0.0) * meshDiffuse;
			vec3 specular = light.specular * pow(max(dot(vs_normal, halfDirection), 0.0), 32) * vec3(1.0, 1.0, 1.0);

			return diffuse + specular;
		}

		vec3 calculatePointLight(pointLight light) {
			vec3 meshDiffuse = vec3(vs_color);
			vec3 viewDirection = normalize(vs_fragPosition - cameraPosition);
			vec3 lightDirection = vs_fragPosition - light.position;
			vec3 halfDirection = -normalize(viewDirection + lightDirection);
			float distance = length(lightDirection);
			float attenuation = max(1 - distance / light.range, 0.0);

			vec3 diffuse = light.diffuse * max(dot(normalize(vs_normal), normalize(-lightDirection)), 0.0) * meshDiffuse;
			vec3 specular = light.specular * pow(max(dot(vs_normal, halfDirection), 0.0), 32) * vec3(1.0, 1.0, 1.0);

			return attenuation * (diffuse + specular);
			return vec3(1.0, 1.0, 1.0);
		}

		void main() {
			vec3 result = ambient * vs_color.rgb;
			for (int i = 0; i < num_dirLight; ++i) {
				result = result + calculateDirLight(dirLights[i]);
			}
			for (int i = 0; i < num_pointLight; ++i) {
				result = result + calculatePointLight(pointLights[i]);
			}
			color = vec4(result, 1.0);
		}
	`,
}

var defaultMeshShader_TEXTURE = glShaderProgram{
	vertexSource: `
		#version 410 core

		layout (location = 0) in vec3 position;
		layout (location = 1) in vec2 uv;
		layout (location = 2) in vec3 normal;

		uniform mat4 model;
		uniform mat4 view;
		uniform mat4 projection;

		out vec2 vs_uv;
		out vec3 vs_normal;
		out vec3 vs_fragPosition;

		void main() {
			vs_uv = uv;
			vs_normal = mat3(transpose(inverse(model))) * normal;
			vs_fragPosition = vec3(model * vec4(position, 1.0));
			gl_Position = projection * view * model * vec4(position, 1.0);
		}
	`,
	fragmentSource: `
		#version 410 core

		struct dirLight {
			vec3 position;
			vec3 direction;
			vec3 diffuse;
			vec3 specular;
		};

		struct pointLight {
			vec3 position;
			float range;
			vec3 diffuse;
			vec3 specular;
		};

		in vec2 vs_uv;
		in vec3 vs_normal;
		in vec3 vs_fragPosition;

		uniform vec3 ambient;

		uniform sampler2D diffuseMap;

		uniform dirLight dirLights[10];
		uniform int num_dirLight = 0;

		uniform pointLight pointLights[50];
		uniform int num_pointLight = 0;

		uniform vec3 cameraPosition;

		out vec4 color;

		vec3 calculateDirLight(dirLight light) {
			vec3 meshDiffuse = texture(diffuseMap, vs_uv).rgb;
			vec3 viewDirection = normalize(vs_fragPosition - cameraPosition);
			vec3 halfDirection = -normalize(viewDirection + light.direction);

			vec3 diffuse = light.diffuse * max(dot(normalize(vs_normal), normalize(-light.direction)), 0.0) * meshDiffuse;
			vec3 specular = light.specular * pow(max(dot(vs_normal, halfDirection), 0.0), 32) * vec3(1.0, 1.0, 1.0);

			return diffuse + specular;
		}

		vec3 calculatePointLight(pointLight light) {
			vec3 meshDiffuse = texture(diffuseMap, vs_uv).rgb;
			vec3 viewDirection = normalize(vs_fragPosition - cameraPosition);
			vec3 lightDirection = vs_fragPosition - light.position;
			vec3 halfDirection = -normalize(viewDirection + lightDirection);
			float distance = length(lightDirection);
			float attenuation = max(1 - distance / light.range, 0.0);

			vec3 diffuse = light.diffuse * max(dot(normalize(vs_normal), normalize(-lightDirection)), 0.0) * meshDiffuse;
			vec3 specular = light.specular * pow(max(dot(vs_normal, halfDirection), 0.0), 32) * vec3(1.0, 1.0, 1.0);

			return attenuation * (diffuse + specular);
			return vec3(1.0, 1.0, 1.0);
		}

		void main() {
			vec3 result = ambient * texture(diffuseMap, vs_uv).rgb;
			for (int i = 0; i < num_dirLight; ++i) {
				result = result + calculateDirLight(dirLights[i]);
			}
			for (int i = 0; i < num_pointLight; ++i) {
				result = result + calculatePointLight(pointLights[i]);
			}
			color = vec4(result, 1.0);
		}
	`,
}

// ----------------------------------------------------------------------------------------------

var defaultMeshShader_COLOR_DEFERRED = glShaderProgram{
	vertexSource: `
		#version 410 core

		layout (location = 0) in vec3 position;
		layout (location = 1) in vec2 uv;
		layout (location = 2) in vec3 normal;

		uniform mat4 model;
		uniform mat3 TImodel;
		uniform mat4 view;
		uniform mat4 projection;
		uniform vec4 color;

		out vec4 vs_color;
		out vec3 vs_normal;
		out vec3 vs_fragPosition;

		void main() {
			vs_color = color;
			vs_normal = TImodel * normal;
			vs_fragPosition = vec3(model * vec4(position, 1.0));
			gl_Position = projection * view * model * vec4(position, 1.0);
		}
	`,
	fragmentSource: `
		#version 410 core

		layout (location = 0) out vec3 gPosition;
		layout (location = 1) out vec3 gNormal;
		layout (location = 2) out vec4 gDiffuse;

		in vec4 vs_color;
		in vec3 vs_normal;
		in vec3 vs_fragPosition;

		uniform float recvShadow = 0.0;

		void main() {
			gPosition = vs_fragPosition;
			gNormal = vs_normal;
			gDiffuse = vec4(vs_color.rgb, recvShadow);
		}
	`,
}

var defaultMeshShader_TEXTURE_DEFERRED = glShaderProgram{
	vertexSource: `
		#version 410 core

		layout (location = 0) in vec3 position;
		layout (location = 1) in vec2 uv;
		layout (location = 2) in vec3 normal;

		uniform mat4 model;
		uniform mat3 TImodel;
		uniform mat4 view;
		uniform mat4 projection;

		out vec2 vs_uv;
		out vec3 vs_normal;
		out vec3 vs_fragPosition;

		void main() {
			vs_uv = uv;
			vs_normal = TImodel * normal;
			vs_fragPosition = vec3(model * vec4(position, 1.0));
			gl_Position = projection * view * model * vec4(position, 1.0);
		}
	`,
	fragmentSource: `
		#version 410 core

		layout (location = 0) out vec3 gPosition;
		layout (location = 1) out vec3 gNormal;
		layout (location = 2) out vec4 gDiffuse;

		in vec2 vs_uv;
		in vec3 vs_normal;
		in vec3 vs_fragPosition;

		uniform sampler2D diffuseMap;

		uniform float recvShadow = 0.0;

		void main() {
			gPosition = vs_fragPosition;
			gNormal = vs_normal;
			gDiffuse = vec4(texture(diffuseMap, vs_uv).rgb, recvShadow);
		}
	`,
}

var defaultDeferredShader_NOLIGHT = glShaderProgram{
	vertexSource: `
		#version 410 core

		layout (location = 0) in vec3 position;
		layout (location = 1) in vec2 uv;

		out vec2 vs_uv;

		void main() {
			vs_uv = uv;
			gl_Position = vec4(position, 1.0);
		}
	`,
	fragmentSource: `
		#version 410 core

		in vec2 vs_uv;

		uniform sampler2D gDiffuse;

		out vec4 color;

		void main() {
			color = vec4(texture(gDiffuse, vs_uv).rgb, 1.0);
		}
	`,
}

var defaultDeferredShader = glShaderProgram{
	vertexSource: `
		#version 410 core

		layout (location = 0) in vec3 position;
		layout (location = 1) in vec2 uv;

		out vec2 vs_uv;

		void main() {
			vs_uv = uv;
			gl_Position = vec4(position, 1.0);
		}
	`,
	fragmentSource: `
		#version 410 core

		struct dirLight {
			vec3 position;
			vec3 direction;
			vec3 diffuse;
			vec3 specular;
		};

		struct pointLight {
			vec3 position;
			float range;
			vec3 diffuse;
			vec3 specular;
		};

		in vec2 vs_uv;

		uniform sampler2D gPosition;
		uniform sampler2D gNormal;
		uniform sampler2D gDiffuse;

		uniform vec3 ambient;

		uniform dirLight dirLights[10];
		uniform int num_dirLight = 0;

		uniform pointLight pointLights[50];
		uniform int num_pointLight = 0;

		uniform vec3 cameraPosition;

		out vec4 color;

		vec3 calculateDirLight(dirLight light) {
			vec3 meshDiffuse = texture(gDiffuse, vs_uv).rgb;
			vec3 vs_fragPosition = texture(gPosition, vs_uv).rgb;
			vec3 vs_normal = texture(gNormal, vs_uv).rgb;

			vec3 viewDirection = normalize(vs_fragPosition - cameraPosition);
			vec3 halfDirection = -normalize(viewDirection + light.direction);

			vec3 diffuse = light.diffuse * max(dot(normalize(vs_normal), normalize(-light.direction)), 0.0) * meshDiffuse;
			vec3 specular = light.specular * pow(max(dot(vs_normal, halfDirection), 0.0), 32) * vec3(1.0, 1.0, 1.0);

			return diffuse + specular;
		}

		vec3 calculatePointLight(pointLight light) {
			vec3 meshDiffuse = texture(gDiffuse, vs_uv).rgb;
			vec3 vs_fragPosition = texture(gPosition, vs_uv).rgb;
			vec3 vs_normal = texture(gNormal, vs_uv).rgb;

			vec3 viewDirection = normalize(vs_fragPosition - cameraPosition);
			vec3 lightDirection = vs_fragPosition - light.position;
			vec3 halfDirection = -normalize(viewDirection + lightDirection);
			float distance = length(lightDirection);
			float attenuation = max(1 - distance / light.range, 0.0);

			vec3 diffuse = light.diffuse * max(dot(normalize(vs_normal), normalize(-lightDirection)), 0.0) * meshDiffuse;
			vec3 specular = light.specular * pow(max(dot(vs_normal, halfDirection), 0.0), 32) * vec3(1.0, 1.0, 1.0);

			return attenuation * (diffuse + specular);
			return vec3(1.0, 1.0, 1.0);
		}

		void main() {
			vec3 result = ambient * texture(gDiffuse, vs_uv).rgb;
			for (int i = 0; i < num_dirLight; ++i) {
				result = result + calculateDirLight(dirLights[i]);
			}
			for (int i = 0; i < num_pointLight; ++i) {
				result = result + calculatePointLight(pointLights[i]);
			}
			color = vec4(result, 1.0);
		}
	`,
}

var defaultShadowMapShader_DIRLIGHT = glShaderProgram{
	vertexSource: `
		#version 410 core

		layout (location = 0) in vec3 position;

		uniform mat4 lightMatrix;
		uniform mat4 model;

		void main() {
			gl_Position = lightMatrix * model * vec4(position, 1.0);
		}
	`,
	fragmentSource: `
		#version 410 core

		void main() {
		}
	`,
}

var defaultBlendShadowShader_DIRLIGHT_DEFERRED = glShaderProgram{
	vertexSource: `
		#version 410 core

		layout (location = 0) in vec3 position;
		layout (location = 1) in vec2 uv;

		out vec2 vs_uv;

		void main() {
			vs_uv = uv;
			gl_Position = vec4(position, 1.0);
		}
	`,
	fragmentSource: `
		#version 410 core

		in vec2 vs_uv;

		uniform vec3 ambient;
		uniform mat4 lightMatrix;
		uniform sampler2D sDirMap;
		uniform sampler2D gPosition;
		uniform sampler2D gDiffuse;

		out vec4 color;

		void main() {
			vec3 fragPos = texture(gPosition, vs_uv).xyz;
			vec3 diffuse = texture(gDiffuse, vs_uv).rgb;
			vec4 fragLightPos = lightMatrix * vec4(fragPos, 1.0);
			vec3 projPos = fragLightPos.xyz / fragLightPos.w;
			projPos = projPos * 0.5 + 0.5;
			float recvShadow = texture(gDiffuse, vs_uv).a;

			float currentDepth = projPos.z;

			float closetDepth = texture(sDirMap, projPos.xy).r;

			float shadow = currentDepth - 0.005 > closetDepth ? 1.0 : 0.0;

			color = vec4(ambient * diffuse, recvShadow * shadow);
		}
	`,
}

// ----------------------------------------------------------------------------------------------

var defaultShadowMapShader_DIRLIGHT_DEFERRED_DEBUG = glShaderProgram{
	vertexSource: `
		#version 410 core

		layout (location = 0) in vec3 position;
		layout (location = 1) in vec2 uv;

		out vec2 vs_uv;

		void main() {
			vs_uv = uv;
			gl_Position = vec4(position, 1.0);
		}
	`,
	fragmentSource: `
		#version 410 core

		in vec2 vs_uv;

		uniform mat4 lightMatrix;

		uniform sampler2D sDirMap;
		uniform sampler2D gPosition;

		out vec4 color;

		void main() {
			vec3 fragPos = texture(gPosition, vs_uv).xyz;
			vec4 fragLightPos = lightMatrix * vec4(fragPos, 1.0);
			vec3 projPos = fragLightPos.xyz / fragLightPos.w;
			projPos = projPos * 0.5 + 0.5;

			color = vec4(vec3(texture(sDirMap, projPos.xy).r), 1.0);
		}
	`,
}
