#include <set>
#include <string>
#include <vector>
using namespace std;

string longestBuildable(vector<string>& words) {
    set<string> have(words.begin(), words.end());
    string best;
    for (auto& w : have) {
        bool ok = true;
        for (size_t p = 1; p < w.size() && ok; p++) ok = have.count(w.substr(0, p)) > 0;
        if (ok && w.size() > best.size()) best = w;
    }
    return best;
}
