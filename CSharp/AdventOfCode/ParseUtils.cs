using SadRogue.Primitives.GridViews;

namespace AdventOfCode;

public static class ParseUtils
{
    public static ArrayView<char> CharGridFrom(string input)
    {
        var width = input.IndexOf(Environment.NewLine, StringComparison.Ordinal);
        return new ArrayView<char>(input.ReplaceLineEndings("").ToCharArray(), width);
    }
}