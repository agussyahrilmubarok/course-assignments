package com.example.witrack.backend.service.v1;

import com.example.witrack.backend.domain.User;
import com.example.witrack.backend.model.AuthDTO;
import com.example.witrack.backend.model.UserDTO;
import com.example.witrack.backend.repos.UserRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.Collections;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class AuthServiceImpl implements AuthService {

    private final UserRepository userRepository;

    @Override
    public AuthDTO.AuthResponse signUp(AuthDTO.SignUpRequest param) {
        User user = new User();
        user.setFullName(param.getFullName());
        user.setEmail(param.getEmail());
        user.setPassword(param.getPassword());
        user.setRoles(Collections.singleton(User.Role.ROLE_USER));
        userRepository.save(user);

        return AuthDTO.AuthResponse.builder()
                .token("")
                .user(UserDTO.UserResponse.fromUser(user))
                .build();
    }

    @Override
    public AuthDTO.AuthResponse signIn(AuthDTO.SignInRequest param) {
        User user = userRepository.findByEmail(param.getEmail())
                .orElseThrow();

        return AuthDTO.AuthResponse.builder()
                .token("")
                .user(UserDTO.UserResponse.fromUser(user))
                .build();
    }

}
