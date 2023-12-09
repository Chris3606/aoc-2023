namespace AdventOfCode;

public readonly record struct NodeData(string Left, string Right);

public sealed class Day08 : BaseDay
{
    private readonly  string _directions;
    private readonly Dictionary<string, NodeData> _nodes = new();

    public Day08()
    {
        var lines = File.ReadAllLines(InputFilePath);

        _directions = lines[0];
        foreach (var line in lines[2..])
            _nodes.Add(line[..3], new NodeData(line[7..10], line[12..15]));
    }

    public override ValueTask<string> Solve_1()
        => new(GetPathLength(_nodes, _directions, "AAA", i => i == "ZZZ").ToString());

    // Trick is that all paths are loops; so we can just take LCM of all path lengths
    public override ValueTask<string> Solve_2()
    {
        var pathLengths = _nodes
            .Where(i => i.Key[2] == 'A')
            .Select(i => GetPathLength(_nodes, _directions, i.Key, n => n[2] == 'Z'))
            .ToArray();
        
        return new(Utility.LCM(pathLengths[0], pathLengths[1], pathLengths[2..]).ToString());
    }

    private static long GetPathLength(Dictionary<string, NodeData> nodes, string directions, string startNode,
                                     Func<string, bool> isEndNode)
    {
        var curElem = startNode;
        int curDirIdx = 0;
        int steps = 0;

        while (!isEndNode(curElem))
        {
            curElem = directions[curDirIdx] switch
            {
                'L' => nodes[curElem].Left,
                'R' => nodes[curElem].Right,
                _ => throw new Exception("Invalid direction")
            };

            steps++;
            curDirIdx = (curDirIdx + 1) % directions.Length;
        }

        return steps;
    }
}
