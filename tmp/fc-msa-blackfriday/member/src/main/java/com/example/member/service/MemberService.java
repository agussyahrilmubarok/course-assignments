package com.example.member.service;

import com.example.member.model.MemberDTO;

public interface MemberService {

    void signUp(MemberDTO.SignUp param);

    MemberDTO.ResponseWithToken signIn(MemberDTO.SignIn param);

    MemberDTO.ResponseWithToken validateToken(String token);
}
