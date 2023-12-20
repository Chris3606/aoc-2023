namespace AdventOfCode;

public readonly record struct SequenceValue(string Label, char Operation, int FocalLen)
{
    public static SequenceValue Parse(string str)
    {
        if (str[^1] == '-')
            return new SequenceValue(str[..^1], '-', -1);

        var split = str.Split('=');
        return new SequenceValue(split[0], '=', int.Parse(split[1]));
        
    }
}

public readonly record struct Lens(string Label, int FocalLength);

public sealed class Day15 : BaseDay
{
    private readonly string[] _input;
    private readonly SequenceValue[] _sequence;

    public Day15()
    {
        _input = File.ReadAllText(InputFilePath).Split(',').ToArray();
        _sequence = _input.Select(SequenceValue.Parse).ToArray();
    }

    public override ValueTask<string> Solve_1() => new(_input.Select(Hash).Sum().ToString());

    public override ValueTask<string> Solve_2()
    {
        var boxes = new List<Lens>[256];
        for (int i = 0; i < 256; i++)
            boxes[i] = new List<Lens>();

        foreach (var seq in _sequence)
        {
            var box = Hash(seq.Label);
            switch (seq.Operation)
            {
                case '-':
                    boxes[box].RemoveAll(i => i.Label == seq.Label);
                    break;
                case '=':
                    int idx = boxes[box].FindIndex(i => i.Label == seq.Label);
                    if (idx == -1)
                        boxes[box].Add(new Lens(seq.Label, seq.FocalLen));
                    else
                        boxes[box][idx] = new Lens(seq.Label, seq.FocalLen);
                    break;
                default:
                    throw new Exception("Unsupported operation");
            }
        }

        int focusingPower = 0;
        for (int i = 0; i < boxes.Length; i++)
        {
            for (int j = 0; j < boxes[i].Count; j++)
                focusingPower += (1 + i) *  (1 + j) * boxes[i][j].FocalLength;
        }

        return new(focusingPower.ToString());
    }

    private static int Hash(string str)
    {
        int curVal = 0;
        foreach (var ch in str)
        {
            curVal += ch;
            curVal *= 17;
            curVal %= 256;
        }

        return curVal;
    }
}
