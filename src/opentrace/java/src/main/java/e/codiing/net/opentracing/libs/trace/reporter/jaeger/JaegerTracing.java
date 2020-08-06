package e.codiing.net.opentracing.libs.trace.reporter.jaeger;

import io.jaegertracing.Configuration;
import io.jaegertracing.Configuration.ReporterConfiguration;
import io.jaegertracing.Configuration.SamplerConfiguration;
import io.jaegertracing.internal.JaegerTracer;
import io.jaegertracing.internal.samplers.ConstSampler;

public final class JaegerTracing {
    private JaegerTracing() {
    }

    public static JaegerTracer init(String service) {
        SamplerConfiguration samplerConfig = SamplerConfiguration.fromEnv()
                .withType(ConstSampler.TYPE)
                .withParam(1)
                .withManagerHostPort("127.0.0.1:6831");

        ReporterConfiguration reporterConfig = ReporterConfiguration.fromEnv()
                .withLogSpans(true);

        Configuration config = new Configuration(service)
                .withSampler(samplerConfig)
                .withReporter(reporterConfig);

        return config.getTracer();
    }

}
