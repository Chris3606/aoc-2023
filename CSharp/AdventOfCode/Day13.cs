using SadRogue.Primitives.GridViews;

namespace AdventOfCode;

public enum ReflectionAxis
{
    Horizontal,
    Vertical
}

public record struct Reflection(ReflectionAxis Axis, int Value);

public sealed class Day13 : BaseDay
{
    private readonly ArrayView<char>[] _patterns;

    public Day13()
    {
        _patterns = File.ReadAllText(InputFilePath)
            .Split(Environment.NewLine + Environment.NewLine)
            .Select(ParseUtils.CharGridFrom)
            .ToArray();
    }

    public override ValueTask<string> Solve_1()
        =>  new (_patterns
            .Select(i => FindReflection(i))
            .Select(i => i ?? throw new Exception("Expected reflection"))
            .Select(ScoreReflection)
            .Sum()
            .ToString());

    public override ValueTask<string> Solve_2()
    {
        int sum = 0;
        foreach (var pattern in _patterns)
        {
            var origReflection = FindReflection(pattern) ?? throw new Exception("No reflection found for initial pattern.");

            foreach (var pos in pattern.Positions())
            {
                char old = pattern[pos];
                pattern[pos] = old == '#' ? '.' : '#';

                var newReflection = FindReflection(pattern, origReflection);
                if (newReflection is not null)
                {
                    sum += ScoreReflection(newReflection.Value);
                    pattern[pos] = old;
                    break;
                }

                pattern[pos] = old;
            }
        }

        return new(sum.ToString());
    }

    private static bool FindHorizontalReflection(IGridView<char> pattern, out int line, int? ignoreLine = null)
    {
        for (int mirrorLine = 1; mirrorLine < pattern.Height; mirrorLine++)
        {
            bool isMirror = true;
            for (int y = mirrorLine; y < pattern.Height; y++)
            {
                int mirrorY = mirrorLine - (y - mirrorLine + 1);
                if (mirrorY < 0)
                    break;
                
                for (int x = 0; x < pattern.Width; x++)
                {
                    if (pattern[x, y] != pattern[x, mirrorY])
                    {
                        isMirror = false;
                        break;
                    }
                }
                
                if (!isMirror)
                    break;
            }

            if (isMirror && mirrorLine != ignoreLine)
            {
                line = mirrorLine;
                return true;
            }
        }
        
        line = -1;
        return false;
    }
    
    private static bool FindVerticalReflection(IGridView<char> pattern, out int line, int? ignoreLine = null)
    {
        for (int mirrorLine = 1; mirrorLine < pattern.Width; mirrorLine++)
        {
            bool isMirror = true;
            for (int x = mirrorLine; x < pattern.Width; x++)
            {
                int mirrorX = mirrorLine - (x - mirrorLine + 1);
                if (mirrorX < 0)
                    break;
                
                for (int y = 0; y < pattern.Height; y++)
                {
                    if (pattern[x, y] != pattern[mirrorX, y])
                    {
                        isMirror = false;
                        break;
                    }
                }
                
                if (!isMirror)
                    break;
            }

            if (isMirror && mirrorLine != ignoreLine)
            {
                line = mirrorLine;
                return true;
            }
        }
        
        line = -1;
        return false;
    }
    
    private static Reflection? FindReflection(IGridView<char> pattern, Reflection? ignoreReflection = null)
    {
         (int? vertIgnore, int? horizIgnore) = (null, null);
        if (ignoreReflection is not null)
        {
            int? val = ignoreReflection.Value.Value;
            (vertIgnore, horizIgnore) = ignoreReflection.Value.Axis switch
            {
                ReflectionAxis.Horizontal => ((int?)null, val),
                ReflectionAxis.Vertical => (val, null),
                _ => throw new ArgumentOutOfRangeException(nameof(ignoreReflection), ignoreReflection, "Invalid reflection axis")
            };
        }
        
        if (FindHorizontalReflection(pattern, out int l1, horizIgnore))
            return new Reflection(ReflectionAxis.Horizontal, l1);
        
        if (FindVerticalReflection(pattern, out int l2, vertIgnore))
            return new Reflection(ReflectionAxis.Vertical, l2);
        
        return null;
    }

    private static int ScoreReflection(Reflection reflection)
        => reflection.Axis switch
        {
            ReflectionAxis.Vertical => reflection.Value,
            ReflectionAxis.Horizontal => 100 * reflection.Value,
            _ => throw new ArgumentOutOfRangeException(nameof(reflection), reflection, "Invalid reflection axis")
        };
}
