namespace AdventOfCode;

public enum HandStrength
{
    HighCard,
    Pair,
    TwoPair,
    ThreeOfAKind,
    FullHouse,
    FourOfAKind,
    FiveOfAKind,
}

public readonly record struct Hand(char[] Cards, int Bid)
{
    public static Hand FromString(string s)
    {
        // Parse from following format (<cards> <bid>). Example: "32T3K 765"
        var cards = s[..5].ToCharArray();
        var bid = int.Parse(s[6..]);

        return new Hand(cards, bid);
    }

    public HandStrength GetStrength(bool useJokers)
    {
        var histDict = Cards.BuildHistogram();
        if (useJokers)
        {
            int jokers = histDict.GetValueOrDefault('J', 0);
            if (jokers == 5)
                return HandStrength.FiveOfAKind;
            
            histDict.Remove('J');
            var max = histDict.MaxBy(i => i.Value);
            histDict[max.Key] += jokers;
        }
        
        var hist = histDict.OrderByDescending(i => i.Value).ToArray();
        
        return hist[0].Value switch
        {
            1 => HandStrength.HighCard,
            2 => hist[1].Value == 2 ? HandStrength.TwoPair : HandStrength.Pair,
            3 => hist[1].Value == 2 ? HandStrength.FullHouse : HandStrength.ThreeOfAKind,
            4 => HandStrength.FourOfAKind,
            5 => HandStrength.FiveOfAKind,
            _ => throw new ArgumentOutOfRangeException()
        };
    }
}

public record HandComparer(bool UseJokers) : IComparer<Hand>
{
    private static readonly Dictionary<char, int> CardValues = new()
    {
        ['2'] = 2,
        ['3'] = 3,
        ['4'] = 4,
        ['5'] = 5,
        ['6'] = 6,
        ['7'] = 7,
        ['8'] = 8,
        ['9'] = 9,
        ['T'] = 10,
        ['J'] = 11,
        ['Q'] = 12,
        ['K'] = 13,
        ['A'] = 14,
    };

    private static readonly Dictionary<char, int> CardValuesJoker = new()
    {
        ['J'] = 1,
        ['2'] = 2,
        ['3'] = 3,
        ['4'] = 4,
        ['5'] = 5,
        ['6'] = 6,
        ['7'] = 7,
        ['8'] = 8,
        ['9'] = 9,
        ['T'] = 10,
        ['Q'] = 12,
        ['K'] = 13,
        ['A'] = 14,
    };
    
    public int Compare(Hand x, Hand y)
    {
        var values = UseJokers ? CardValuesJoker : CardValues;
        
        var xStrength = x.GetStrength(UseJokers);
        var yStrength = y.GetStrength(UseJokers);

        if (xStrength != yStrength)
            return xStrength.CompareTo(yStrength);

        for (int i = 0; i < x.Cards.Length; i++)
            if (values[x.Cards[i]] != values[y.Cards[i]])
                return values[x.Cards[i]].CompareTo(values[y.Cards[i]]);

        return 0;
    }
}

public sealed class Day07 : BaseDay
{
    private readonly Hand[] _cards;

    public Day07()
    {
        _cards = File.ReadLines(InputFilePath).Select(Hand.FromString).ToArray();
    }

    public override ValueTask<string> Solve_1()
    {
        int answer = _cards.Order(new HandComparer(false)).Select((h, i) => h.Bid * (i + 1)).Sum();
        return new(answer.ToString());
    }
    
    public override ValueTask<string> Solve_2()
    {
        int answer = _cards.Order(new HandComparer(true)).Select((h, i) => h.Bid * (i + 1)).Sum();
        return new(answer.ToString());
    }
}
