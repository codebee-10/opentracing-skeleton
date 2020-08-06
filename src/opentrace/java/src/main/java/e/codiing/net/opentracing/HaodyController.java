package e.codiing.net.opentracing;

import com.google.common.collect.ImmutableMap;
import e.codiing.net.opentracing.libs.trace.reporter.jaeger.JaegerTracing;
import io.jaegertracing.internal.JaegerTracer;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.RequestMapping;
import io.opentracing.Span;

@RestController
public class HaodyController {

//    private final Tracer tracer;
//    private HaodyController(Tracer tracer) {
//        this.tracer = tracer;
//    }

    @RequestMapping("/jwt-auth")
    public String index() {
        String serviceName = "spring_boot UserList";
        try (JaegerTracer tracer = JaegerTracing.init(serviceName)) {
            System.out.println(tracer);
            Span span = tracer.buildSpan("someWork").ignoreActiveSpan().start();
            span.setTag("service name", serviceName);
            String helloStr = String.format("user name, %s!", "ethan");
            span.log(ImmutableMap.of("event", "string-format", "value", helloStr));
            span.finish();
        }
        return "haody, auth successful!";
    }
}
