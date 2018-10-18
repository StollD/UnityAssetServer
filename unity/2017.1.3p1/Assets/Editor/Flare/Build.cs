using UnityEngine;
using UnityEditor;

public class FlareAsset
{
    static void Build()
    {
        // Find the lensflare and it's texture. They will always be named "sun_flare.flare" and "sun_flare.png"
        Object flare = AssetDatabase.LoadAssetAtPath("Assets/Editor/Flare/Flare.flare", typeof(Flare));
        Object tex = AssetDatabase.LoadAssetAtPath("Assets/Editor/Flare/Flare.png", typeof(Texture2D));
        
        // Build the bundle
        BuildPipeline.BuildAssetBundle(flare, new[] {flare, tex}, "Assets/Editor/Flare/Flare.unity3d", BuildAssetBundleOptions.ChunkBasedCompression, BuildTarget.StandaloneWindows);
    }
}