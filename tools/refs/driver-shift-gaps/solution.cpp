#include <map>
#include <vector>
using namespace std;

vector<vector<int>> understaffed(vector<vector<int>>& shifts, int m, int dayEnd) {
    map<int, int> change{{0, 0}, {dayEnd, 0}};
    for (auto& s : shifts) change[s[0]]++, change[s[1]]--;
    vector<vector<int>> out;
    int on = 0;
    for (auto it = change.begin(); it != change.end() && it->first < dayEnd; ++it) {
        on += it->second;
        int next = std::next(it)->first;
        if (on >= m) continue;
        if (!out.empty() && out.back()[1] == it->first) out.back()[1] = next;
        else out.push_back({it->first, next});
    }
    return out;
}
