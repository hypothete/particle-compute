#version 430 core
#define NUMPARTICLES 1024
layout (local_size_x = 1024, local_size_y = 1) in;
layout (std140, binding = 0) buffer Pos {
  vec4 positions[];
};

layout (std140, binding = 1) buffer Vel {
  vec4 velocities[];
};

void main() {
  uint index = gl_GlobalInvocationID.x + gl_GlobalInvocationID.y * gl_NumWorkGroups.x * gl_WorkGroupSize.x;

	if(index > NUMPARTICLES) {
    return;
  }

  float t = 0.01;
  
  vec3 pPos = positions[index].xyz;
  vec3 vPos = velocities[index].xyz;

  vec3 g = normalize(pPos) * -9.0; // grav strength
  vec3 pp = pPos + vPos * t + 0.5 * t * t * g;
  vec3 vp = vPos + g * t;

  positions[index] = vec4(pp, 1.0);
  velocities[index] = vec4(vp, 0.0);
}