# Add project specific ProGuard rules here.
# Keep Kotlin serialization
-keepattributes *Annotation*, InnerClasses
-dontnote kotlinx.serialization.AnnotationsKt

-keepclassmembers class kotlinx.serialization.json.** {
    *** Companion;
}
-keepclasseswithmembers class kotlinx.serialization.json.** {
    kotlinx.serialization.KSerializer serializer(...);
}

# Keep EcoMessage classes for serialization
-keep class dev.eco.core.protocol.** { *; }
-keepclassmembers class dev.eco.core.protocol.** { *; }
