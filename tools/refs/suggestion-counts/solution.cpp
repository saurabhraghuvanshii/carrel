#include <string>
#include <vector>
using namespace std;

vector<int> countMatches(vector<string>& words, vector<string>& prefixes) {
    vector<int> out;
    for (auto& p : prefixes) {
        int count = 0;
        for (auto& w : words)
            if (w.compare(0, p.size(), p) == 0) count++;
        out.push_back(count);
    }
    return out;
}
