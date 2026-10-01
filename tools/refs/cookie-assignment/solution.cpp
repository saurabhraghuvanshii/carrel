#include <algorithm>
#include <vector>
using namespace std;

int feedChildren(vector<int>& greed, vector<int>& cookies) {
    sort(greed.begin(), greed.end());
    sort(cookies.begin(), cookies.end());
    size_t fed = 0;
    for (int size : cookies)
        if (fed < greed.size() && size >= greed[fed]) fed++;
    return fed;
}
