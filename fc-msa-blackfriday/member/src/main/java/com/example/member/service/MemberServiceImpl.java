package com.example.member.service;

import com.auth0.jwt.interfaces.DecodedJWT;
import com.example.member.domain.User;
import com.example.member.exception.InvalidPasswordException;
import com.example.member.exception.UserAlreadyExistsException;
import com.example.member.exception.UserNotFoundException;
import com.example.member.model.MemberDTO;
import com.example.member.repos.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
@Slf4j
@RequiredArgsConstructor
public class MemberServiceImpl implements MemberService {

    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;
    private final JWTService jwtService;

    @Override
    public void signUp(MemberDTO.SignUp param) {
        boolean exists = userRepository.findByUsername(param.getUsername()).isPresent();
        if (exists) {
            log.warn("Sign-up failed: username '{}' already exists", param.getUsername());
            throw new UserAlreadyExistsException("Username is already used");
        }

        User user = new User();
        user.setId(UUID.randomUUID().toString());
        user.setUsername(param.getUsername());
        user.setPassword(passwordEncoder.encode(param.getPassword()));

        userRepository.save(user);
        log.info("User '{}' registered successfully with ID={}", user.getUsername(), user.getId());
    }

    @Override
    public MemberDTO.ResponseWithToken signIn(MemberDTO.SignIn param) {
        User user = userRepository.findByUsername(param.getUsername())
                .orElseThrow(() -> {
                    log.warn("Sign-in failed: user '{}' not found", param.getUsername());
                    return new UserNotFoundException("User not found: " + param.getUsername());
                });

        if (!passwordEncoder.matches(param.getPassword(), user.getPassword())) {
            log.warn("Sign-in failed: invalid password for username='{}'", param.getUsername());
            throw new InvalidPasswordException("Invalid password");
        }

        String token = jwtService.generate(user);
        log.info("User '{}' signed in successfully, token generated", param.getUsername());

        return MemberDTO.ResponseWithToken.builder()
                .token(token)
                .member(MemberDTO.Response.from(user))
                .build();
    }

    @Override
    public MemberDTO.ResponseWithToken validateToken(String token) {
        DecodedJWT claims = jwtService.validateToken(token);
        String username = claims.getSubject();

        log.debug("Token subject extracted: username={}", username);

        User user = userRepository.findByUsername(username)
                .orElseThrow(() -> {
                    log.error("Token validation failed: user not found with username={}", username);
                    return new UserNotFoundException("User not found with username " + username);
                });

        log.info("Token validated successfully for username={}", username);

        return MemberDTO.ResponseWithToken.builder()
                .token(token)
                .member(MemberDTO.Response.from(user))
                .build();
    }
}
