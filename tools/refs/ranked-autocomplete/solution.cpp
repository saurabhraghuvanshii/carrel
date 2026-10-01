#include <algorithm>
#include <map>
#include <string>
#include <vector>
using namespace std;

// Checks every stored word; simple, and fine at this size.
class Autocomplete {
    map<string, int> total;

public:
    Autocomplete() {}

    void add(string word, int score) { total[word] += score; }

    vector<string> suggest(string prefix) {
        vector<pair<int, string>> match;  // {-total, word}
        for (auto& [word, t] : total)
            if (word.compare(0, prefix.size(), prefix) == 0) match.push_back({-t, word});
        sort(match.begin(), match.end());
        vector<string> out;
        for (size_t i = 0; i < match.size() && i < 3; i++) out.push_back(match[i].second);
        return out;
    }
};
