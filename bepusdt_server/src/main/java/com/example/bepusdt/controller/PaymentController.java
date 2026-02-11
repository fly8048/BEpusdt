package com.example.bepusdt.controller;

import com.example.bepusdt.dto.CancelOrderRequest;
import com.example.bepusdt.dto.CreateOrderRequest;
import com.example.bepusdt.service.BEpusdtService;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
@RequestMapping("/api/payment")
public class PaymentController {

    private final BEpusdtService bepusdtService;

    public PaymentController(BEpusdtService bepusdtService) {
        this.bepusdtService = bepusdtService;
    }

    /**
     * 创建订单
     * 客户端 -> Spring Boot -> BEpusdt
     */
    @PostMapping("/create")
    public String createOrder(@RequestBody CreateOrderRequest request) {
        // 👉 这里你后面可以加：
        // 1. 登录校验
        // 2. 订单是否重复
        // 3. 金额校验
        // 4. 风控逻辑

        return bepusdtService.createOrder(request);
    }

    /**
     * 取消订单
     */
    @PostMapping("/cancel")
    public String cancelOrder(@RequestBody CancelOrderRequest request) {
        return bepusdtService.cancelOrder(request);
    }

//    /**
//     * 支付回调（BEpusdt -> 你的系统）
//     */
//    @PostMapping("/notify")
//    public String notify(@RequestBody Map<String, Object> notifyData) {
//        System.out.println("====== BEpusdt 支付回调 ======");
//        notifyData.forEach((k, v) -> System.out.println(k + " = " + v));
//
//        // 👉 后面你可以在这里：
//        // 1. 验签
//        // 2. 更新订单状态
//        // 3. 做幂等处理
//        // 4. 发 MQ / Webhook / WebSocket
//
//        return "success";
//    }
}


//@RestController
//@RequestMapping("/api/payment")
//public class PaymentController {
//
//    private final BEpusdtService service;
//
//    public PaymentController(BEpusdtService service) {
//        this.service = service;
//    }
//
//    @PostMapping("/create")
//    public String create(@RequestBody CreateOrderRequest request) {
//        return service.createOrder(request);
//    }
//
//    @PostMapping("/cancel")
//    public String cancel(@RequestBody CancelOrderRequest request) {
//        return service.cancelOrder(request);
//    }
//
//    @PostMapping("/notify")
//    public String notify(@RequestBody Map<String, Object> body) {
//        System.out.println("====== BEpusdt 回调 ======");
//        body.forEach((k, v) -> System.out.println(k + " = " + v));
//        return "success";
//    }
//}


