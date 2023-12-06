namespace AdventOfCode;

public sealed class Day01 : BaseDay
{
    private static readonly Dictionary<string, char> WordToDigitMap = new()
    {
        ["one"] = '1',
        ["two"] = '2',
        ["three"] = '3',
        ["four"] = '4',
        ["five"] = '5',
        ["six"] = '6',
        ["seven"] = '7',
        ["eight"] = '8',
        ["nine"] = '9',
        ["zero"] = '0'
    };
    
    private readonly string[] _input;

    public Day01()
    {
        _input = File.ReadAllLines(InputFilePath);
    }
    
    private char FindFirstDigit(string s, bool words = false)
    {
        for (int i = 1; i <= s.Length; i++)
        {
            if (char.IsDigit(s[i - 1]))
                return s[i - 1];
            
            if (!words) continue;
            var match = WordToDigitMap.FirstOrDefault(pair => s.AsSpan(..i).EndsWith(pair.Key)).Value;
            if (match != default) return match;
        }

        throw new Exception("No digit found");
    }
    
    private char FindLastDigit(string s, bool words = false)
    {
        for (int i = s.Length - 1; i >= 0; i--)
        {
            if (char.IsDigit(s[i]))
                return s[i];
        
            if (!words) continue;
            var match = WordToDigitMap.FirstOrDefault(pair => s.AsSpan(i..).StartsWith(pair.Key)).Value;
            if (match != default) return match;
        }

        throw new Exception("No digit found");
    }

    private int ParseDigit(string line)
        => int.Parse($"{FindFirstDigit(line)}{FindLastDigit(line)}");
    
    private int ParseDigitWithWords(string line)
        => int.Parse($"{FindFirstDigit(line, true)}{FindLastDigit(line, true)}");

    public override ValueTask<string> Solve_1() => new(_input.Select(ParseDigit).Sum().ToString());

    public override ValueTask<string> Solve_2() => new(_input.Select(ParseDigitWithWords).Sum().ToString());
}
