package com.example.bepusdt.util;

import java.security.MessageDigest;
import java.util.*;

public class SignUtil {

    public static String sign(Map<String, Object> params, String apiToken) {
        List<String> keys = new ArrayList<>();

        for (Map.Entry<String, Object> entry : params.entrySet()) {
            if (entry.getValue() != null
                    && !"".equals(entry.getValue())
                    && !"signature".equals(entry.getKey())) {
                keys.add(entry.getKey());
            }
        }

        Collections.sort(keys);

        StringBuilder sb = new StringBuilder();
        for (String key : keys) {
            sb.append(key).append("=")
                    .append(params.get(key))
                    .append("&");
        }

        // 去掉最后一个 &
        sb.deleteCharAt(sb.length() - 1);
        sb.append(apiToken);

        return md5(sb.toString());
    }

    private static String md5(String data) {
        try {
            MessageDigest md = MessageDigest.getInstance("MD5");
            byte[] digest = md.digest(data.getBytes());
            StringBuilder hex = new StringBuilder();
            for (byte b : digest) {
                String s = Integer.toHexString(0xff & b);
                if (s.length() == 1) hex.append("0");
                hex.append(s);
            }
            return hex.toString();
        } catch (Exception e) {
            throw new RuntimeException("MD5 error", e);
        }
    }
}

